package route

import (
	"github.com/vtrifonov/http-api-mock/definition"
)

//Router contains the functions to check the http request and return the matching mock.
type Router interface {
	Route(req *definition.Request) (*definition.Mock, map[string]string)
	SetMockDefinitions(mocks []definition.Mock)
}
