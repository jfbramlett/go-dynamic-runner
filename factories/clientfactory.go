package factories

import (
	"fmt"
	"reflect"
	"strings"
)


// ClientFactory is an interface used to maintain the set of Grpc clients that are used for our tests
type ClientFactory interface {
	GetClient(name string) (interface{}, error)
	RegisterClient(client interface{})
}

type defaultClientFactory struct {
	registeredClients		map[string]interface{}
}

// GetClient retrieves the appropriate client given based on the specified name
func (r *defaultClientFactory) GetClient(name string) (interface{}, error) {
	client, found := r.registeredClients[name]
	if !found {
		err := fmt.Errorf("failed to find registered client with name %s", name)
		return nil, err
	}

	return client, nil
}

// RegisterClient registers a given client under a given name. The name is defined as the name of the type of the given
// client with the last element starting with a capital letter. For example, if the type name is routeguide.routeDetails
// the registered name will be routeguide.RouteDetails (this follows the normal Grpc naming pattern)
func (r *defaultClientFactory) RegisterClient(client interface{}) {
	r.registeredClients[r.getName(client)] = client
}


func (r *defaultClientFactory) getName(client interface{}) string {
	typeOf := reflect.TypeOf(client)
	typeName := typeOf.String()
	idx := strings.LastIndex(typeName, ".")
	if idx > -1 {
		base := typeName[:idx]
		structName := typeName[idx+1:]
		typeName = fmt.Sprintf("%s.%s", base, strings.Title(structName))
	}

	if isPointer(typeOf) {
		typeName = typeName[1:]
	}

	return typeName
}

// GlobalClientFactory is a global variable used to hold instances of our registered Grpc clients
var GlobalClientFactory = &defaultClientFactory{registeredClients: make(map[string]interface{})}