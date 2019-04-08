package factories

import (
	"fmt"
)

type ClientFactory interface {
	GetClient(name string) (interface{}, error)
	RegisterClient(name string, client interface{})
}

type defaultClientFactory struct {
	registeredClients		map[string]interface{}
}


func (r *defaultClientFactory) GetClient(name string) (interface{}, error) {
	client, found := r.registeredClients[name]
	if !found {
		err := fmt.Errorf("failed to find registered client with name %s", name)
		return nil, err
	}

	return client, nil
}

func (r *defaultClientFactory) RegisterClient(name string, client interface{}) {
	r.registeredClients[name] = client
}


var GlobalClientFactory = &defaultClientFactory{registeredClients: make(map[string]interface{})}