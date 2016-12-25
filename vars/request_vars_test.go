package vars

import (
	"testing"

	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/persist"
	"github.com/vtrifonov/http-api-mock/vars/fakedata"
)

func getVarsProcessor() VarsProcessor {
	filePersist := persist.NewFilePersister("./test_persist")
	persistBag := persist.GetNewPersistEngineBag(filePersist)
	return VarsProcessor{FillerFactory: MockFillerFactory{}, FakeAdapter: fakedata.NewDummyDataFaker("AleixMG"), PersistEngines: persistBag}
}

func TestRequestVars_RouteParameterPath(t *testing.T) {
	processor := getVarsProcessor()

	req := &definition.Request{}
	req.Path = "/users/1"

	mock := &definition.Mock{}
	mock.Request.Path = "/users/:userId"
	mock.Response.Body = "{{ request.path.userId }}"

	processor.Eval(req, mock)

	if mock.Response.Body != "1" {
		t.Error("The response body should have the value of 1", mock.Response.Body)
	}
}

func TestRequestVars_MultipleRouteParametersPath(t *testing.T) {
	processor := getVarsProcessor()

	req := &definition.Request{}
	req.Path = "/users/administrators/1"

	mock := &definition.Mock{}
	mock.Request.Path = "/users/:role/:userId"
	mock.Response.Body = "{{ request.path.role }}/{{ request.path.userId }}"

	processor.Eval(req, mock)

	if mock.Response.Body != "administrators/1" {
		t.Error("The response body should have the value of administrators/", mock.Response.Body)
	}
}

func TestRequestVars_GlobPath(t *testing.T) {
	processor := getVarsProcessor()

	req := &definition.Request{}
	req.Path = "/users/1"

	mock := &definition.Mock{}
	mock.Request.Path = "/users/*"
	mock.Response.Body = "{{ request.url./users/(?P<value>\\d+) }}"

	processor.Eval(req, mock)

	if mock.Response.Body != "1" {
		t.Error("The response body should have the value of 1", mock.Response.Body)
	}
}
