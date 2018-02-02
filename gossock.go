//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"reflect"
	"context"
)

//
// Gossock is Nossock protocol adapter
// that wraps io.ReadWriteCloser
//
type Gossock struct {

	// conn is underlying connection
	conn io.ReadWriteCloser

	// handlers is user callbacks for messages
	handlers map[string][]interface{}

	// registry is storage of message types
	// to message name mapping and vice-versa
	registry *Registry

	// parser does parsing from conn reader
	parser *parser

	// serializer does serialization to conn writer
	serializer *serializer

	// Close is channel to handle Gossock close
	// due to error or because of any other case
	Errors chan error

	// context to pass to handlers
	ctx context.Context
}

//
// New creates new Gossock around io.ReadWriteCloser
//
func New(registry *Registry) *Gossock {

	g := &Gossock{
		handlers:   make(map[string][]interface{}),
		registry:   registry,
		Errors:     make(chan error, 1),
	}

	return g
}

//
// Starts loop
//
func (g *Gossock) Start(ctx context.Context, conn io.ReadWriteCloser) {

	g.ctx = ctx
	g.conn = conn
	g.parser = newParser(conn)
	g.serializer = newSerializer(conn)

	go g.loop()
}

//
// Closes underlying io.ReadWriteCloser
//
func (g *Gossock) Close() error {
	return g.conn.Close()
}

//
// loops represents inner parsing / frames processing loop
//
func (g *Gossock) loop() {

	for {
		var frame frame
		var ok bool

		select {
		case err := <-g.parser.errors:
			g.Errors <- err
			close(g.Errors)
			return

		case frame = <-g.parser.frames:
		}

		var typ reflect.Type = nil

		if typ, ok = g.registry.registryInverse[frame.name]; !ok {
			// frame is not registered
			continue
		}

		var obj interface{} = nil

		switch frame.typ {
		case frameTypeJson:
			var err error
			obj, err = g.parseBody(typ, frame.body)
			if err != nil {
				continue
			}
		case frameTypeBinary:
			obj = &obj
			obj = reflect.ValueOf(&frame.body).Convert(reflect.PtrTo(typ)).Interface()
		default:
			// frame of unknown type
			continue
		}

		/* Call registered handlers */
		// TODO: Make optionally async and run in separate goroutine
		for _, handler := range g.handlers[frame.name] {

			func(handler interface{}, obj interface{}) {

				defer func() {
					if r := recover(); r != nil {
						log.Println("Panic recovered", r)
						// Recovered, log it
					}
				}()

				reflect.ValueOf(handler).Call([]reflect.Value{
					reflect.ValueOf(g.ctx),
					reflect.ValueOf(obj),
				})

			}(handler, obj)
		}
	}

}

//
// Send send frame to other party
//
func (g *Gossock) Send(body interface{}) error {

	name, ok := g.registry.registry[reflect.TypeOf(body)]
	if !ok {
		return errors.New("Message for type " + reflect.TypeOf(body).String() + " is not registered")
	}

	var f frame

	// Special case for binary data
	if reflect.TypeOf(body).ConvertibleTo(reflect.TypeOf([]byte{})) {

		f = frame{
			name: name,
			typ:  frameTypeBinary,
			body: reflect.ValueOf(body).Convert(reflect.TypeOf([]byte{})).Interface().([]byte),
		}

	} else {

		var err error
		b, err := g.serializeBody(body)

		if err != nil {
			return err
		}

		f = frame{
			name: name,
			typ:  frameTypeJson,
			body: b,
		}
	}

	return g.serializer.serialize(&f)
}

//
// On registers new callback handler for frame
//
func (g *Gossock) On(handler interface{}) error {

	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() != 2 {
		panic(errors.New("Handlers should accept exactly two argunents"))
		return errors.New("Handler should accept exactly two arguments")
	}

	if handlerType.In(0).Kind() != reflect.Interface || !reflect.TypeOf((*context.Context)(nil)).Elem().Implements( handlerType.In(0)){
		panic(errors.New("First argument of handler should be context.Context"))
	}

	if handlerType.In(1).Kind() != reflect.Ptr {
		return errors.New("Second argument of handler should accept pointer")
	}

	parameterType := handlerType.In(1).Elem()
	name, ok := g.registry.registry[parameterType]
	if !ok {
		return errors.New("Message not registered")
	}

	_, ok = g.handlers[name]
	if !ok {
		g.handlers[name] = []interface{}{}
	}

	g.handlers[name] = append(g.handlers[name], handler)

	return nil
}

//
// Off unregisters concrete handler for name
// TODO: DO LOCK
//
func (g *Gossock) Off(name string, handler interface{}) {
	handlers := g.handlers[name]

	var na []interface{}
	for _, v := range handlers {
		if v == handler {
			continue
		} else {
			na = append(na, v)
		}
	}

	g.handlers[name] = na
}

//
// OffAll unregisters all handler for name
//
func (g *Gossock) OffAll(name string) {
	delete(g.handlers, name)
}

//
// parseBody parses body to structure of type typ
//
func (g *Gossock) parseBody(typ reflect.Type, body []byte) (interface{}, error) {
	obj := reflect.New(typ).Interface()
	err := json.Unmarshal(body, obj)
	return obj, err
}

//
// serializeBody serializes body to byte array
//
func (g *Gossock) serializeBody(body interface{}) ([]byte, error) {
	return json.Marshal(body)
}
