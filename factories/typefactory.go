package factories

import (
	"fmt"
	"reflect"
	"strings"
)

// TypeFactory represents a factory used to create types in the system (like method parameters)
type TypeFactory interface {
	// GetInstanceCreator finds the instance creator associated with the given name
	GetInstanceCreator(name string) (InstanceCreator, error)

	// GetInstanceCreatorForType finds the appropriate instance create for the given type
	GetInstanceCreatorForType(typ reflect.Type) (InstanceCreator, error)

	// RegisterType registers a new instance creator under the given type name
	RegisterType(typeName string, instanceCreator InstanceCreator)
}

func getDefaultTypeFactory() TypeFactory {
	typeFactory := &defaultTypeFactory{}
	typeFactory.registerTypes()
	return typeFactory
}

func GetTypeName(typ reflect.Type) string {
	typeName := fmt.Sprintf("%s", typ)
	if strings.HasPrefix(typeName, "*") {
		typeName = typeName[1:]
	}
	return typeName
}

type defaultTypeFactory struct {
	typeMap 		map[string]InstanceCreator
}


func (g *defaultTypeFactory) GetInstanceCreator(name string) (InstanceCreator, error) {
	t, found := g.typeMap[name]

	if !found {
		return nil, fmt.Errorf("no type named %s registered", name)
	} else {
		return t, nil
	}
}

func (g *defaultTypeFactory) GetInstanceCreatorForType(typ reflect.Type) (InstanceCreator, error) {
	typeName := typ.String()
	if isPointer(typ) {
		typeName = typeName[1:]
	}

	typeCreator, err := g.GetInstanceCreator(typeName)
	if err != nil {
		typeCreator = NewReflectionInstanceCreator(typ)
		g.typeMap[typeName] = typeCreator
	}

	return typeCreator, nil
}

func (g *defaultTypeFactory) registerTypes() {
	g.typeMap = make(map[string]InstanceCreator)

	g.typeMap["context.Context"] = newContextInstanceCreator()
}


func (g *defaultTypeFactory) RegisterType(typeName string, instanceCreator InstanceCreator) {
	g.typeMap[typeName] = instanceCreator
}

// GlobalTypeFactory is a global variable containing our set of registered types
var GlobalTypeFactory = getDefaultTypeFactory()