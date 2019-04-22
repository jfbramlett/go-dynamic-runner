package factories

import (
	"context"
	"reflect"
)

// InstanceCreator is an interface to a "factory" used to create instances of a given type
type InstanceCreator interface {
	NewInstance() interface{}
}

type reflectionInstanceCreator struct {
	typeOf		reflect.Type
}

// NewInstance creates a new instance of the "type" defined by this instance creator
func (i *reflectionInstanceCreator) NewInstance() interface{} {
	if isPointer(i.typeOf) {
		return reflect.New(i.typeOf.Elem()).Interface()
	} else {
		return reflect.New(i.typeOf).Interface()
	}
}

// NewReflectionInstanceCreator creates a new InstanceCreator for the given type
func NewReflectionInstanceCreator(typeOf reflect.Type) InstanceCreator {
	return &reflectionInstanceCreator{typeOf: typeOf}
}

type contextInstanceCreator struct {}
func (c *contextInstanceCreator) NewInstance() interface{} {
	return context.Background()
}

func newContextInstanceCreator() InstanceCreator {
	return &contextInstanceCreator{}
}


func isPointer(typ reflect.Type) bool {
	return typ.Kind() == reflect.Ptr
}