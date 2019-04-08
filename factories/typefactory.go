package factories

import (
	"fmt"
	"reflect"
	"strings"
)

type TypeFactory interface {
	GetInstanceCreator(name string) (InstanceCreator, error)
	GetInstanceCreatorForType(typ reflect.Type) (InstanceCreator, error)
	RegisterType(typeName string, instanceCreator InstanceCreator)
	Close()
}

func getDefaultTypeFactory() TypeFactory {
	typeFactory := &defaultTypeFactory{}
	typeFactory.registerTypes()
	return typeFactory
}

func GetTypeNameFromIns(ins interface{}) string {
	return GetTypeName(reflect.TypeOf(ins))
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
	if strings.HasPrefix(typeName, "*") {
		typeName = typeName[1:]
	}
	return g.GetInstanceCreator(typeName)
}

func (g *defaultTypeFactory) Close() {}

func (g *defaultTypeFactory) registerTypes() {
	g.typeMap = make(map[string]InstanceCreator)

	g.typeMap["context.Context"] = newContextInstanceCreator()
}


func (g *defaultTypeFactory) RegisterType(typeName string, instanceCreator InstanceCreator) {
	g.typeMap[typeName] = instanceCreator
}


var GlobalTypeFactory = getDefaultTypeFactory()