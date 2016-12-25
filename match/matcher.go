package match

import (
	"github.com/vtrifonov/http-api-mock/definition"
)

//Matcher checks if the received request matches with some specific mock request definition.
type Matcher interface {
	Match(req *definition.Request, mock *definition.Request) (bool, error)
}
