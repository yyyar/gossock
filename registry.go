//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

import (
	"reflect"
	"errors"
)

//
// Registry is storage of all messages
// and their types understood by Gossock
//
type Registry struct {
	registry map[reflect.Type]string
	registryInverse map[string]reflect.Type
}

//
// NewRegistry creates new instance of Registry
//
func NewRegistry() *Registry {

	return &Registry{
		registry: make(map[reflect.Type]string),
		registryInverse:  make(map[string]reflect.Type),
	}
}

//
// Register registers message with name and associates it with type
//
func (r *Registry) Register(name string, typ interface{}) error {

	_, ok := r.registry[reflect.TypeOf(typ)]
	if ok {
		return errors.New(reflect.TypeOf(typ).Name() + " is already registered")
	}

	_, ok = r.registryInverse[name]
	if ok {
		return errors.New(name + " is already registered")
	}

	r.registry[reflect.TypeOf(typ)] = name
	r.registryInverse[name] = reflect.TypeOf(typ)

	return nil
}

