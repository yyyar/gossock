//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

import (
	"encoding/binary"
	"io"
)

//
// parser Represents nossock protocol frames parser
//
type parser struct {

	// reader to read data for parsing
	reader io.Reader

	// frame may be incomplete parsed frame
	frame frame

	// frames is used to push parsed frames
	frames chan (frame)

	errors chan (error)
}

//
// newParser creates new instance of parser
//
func newParser(reader io.Reader) *parser {

	parser := parser{
		reader: reader,
		frames: make(chan frame),
		frame:  frame{},
		errors: make(chan error),
	}

	go parser.loop()

	return &parser
}

//
// loop parses frames until it gets EOF error
// from underlying reader
//
func (p *parser) loop() {

	for {

		if err := p.nextFrame(); err != nil {
			p.errors <- err
			return
		}

		p.frames <- p.frame
	}
}

//
// nextFrame reads bytes from reader, parses it
// and writes data to frame.
//
func (p *parser) nextFrame() error {

	// frame type (1 byte)
	if err := binary.Read(p.reader, binary.BigEndian, &p.frame.typ); err != nil {
		return err
	}

	// nameLen (1 byte)
	var nameLen byte = 0
	if err := binary.Read(p.reader, binary.BigEndian, &nameLen); err != nil {
		return err
	}

	// name (nameLen bytes)
	var name []byte = make([]byte, nameLen)
	if err := binary.Read(p.reader, binary.BigEndian, &name); err != nil {
		return err
	}
	p.frame.name = string(name)

	// bodyLen (4 bytes)
	var bodyLen uint32
	if err := binary.Read(p.reader, binary.BigEndian, &bodyLen); err != nil {
		return err
	}

	// body (bodyLen bytes)
	var body []byte = make([]byte, bodyLen)
	if err := binary.Read(p.reader, binary.BigEndian, &body); err != nil {
		return err
	}
	p.frame.body = body

	return nil
}
