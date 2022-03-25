package elastic

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

type AccessImpl struct {
	es *elasticsearch.Client
}

func NewAccessImpl(suppliedClient *elasticsearch.Client) (*AccessImpl, error) {
	if suppliedClient == nil {
		return nil, fmt.Errorf("cannot create interface with empty client")
	}
	return &AccessImpl{es: suppliedClient}, nil
}

func (a *AccessImpl) CreateRole(name string, roleRequest CreateRoleRequest) (bool, error) {
	jsonRole, err := roleRequest.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonRole)
	res, err := a.es.Security.PutRole(name, buffer)
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	return true, nil
}

func (a *AccessImpl) CreateUser(name string, userCreation CreateUserRequest) (bool, error) {
	jsonUser, err := userCreation.Marshal()
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(jsonUser)
	res, err := a.es.Security.PutUser(name, buffer)
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, fmt.Errorf("%s: %s", errorMessage, res.Status())
	}

	return true, nil
}
