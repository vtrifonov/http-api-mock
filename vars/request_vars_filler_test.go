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

func TestRequestVarsFiller_RouteParameterPath(t *testing.T) {
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

func TestRequestVarsFiller_Test(t *testing.T) {
	processor := getVarsProcessor()

	req := &definition.Request{}
	req.Path = "{ \"test\" : 1 }"

	mock := &definition.Mock{}
	mock.Request.Path = "{ \"test\" : :testValue }"
	mock.Response.Body = "{{ request.path.testValue }}"

	processor.Eval(req, mock)

	if mock.Response.Body != "1" {
		t.Error("The response body should have the value of 1", mock.Response.Body)
	}
}

func TestRequestVarsFiller_MultipleRouteParametersPath(t *testing.T) {
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

func TestRequestVarsFiller_GlobPath(t *testing.T) {
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

func TestRequestVarsFiller_BodyPart(t *testing.T) {
	processor := getVarsProcessor()

	req := &definition.Request{}

	req.Path = "/users/1"
	req.Body = "{ \"a\": { \"aa\": \"nameValue\", \"ab\": 2}, \"b\": 3 }"

	mock := &definition.Mock{}
	mock.Request.Path = "/users/*"
	mock.Response.Body = "{ \"name\": \"{{ request.body.a.aa }}\", \"value\" : {{ request.body.a.ab }} }"

	expectedResult := "{ \"name\": \"nameValue\", \"value\" : 2 }"

	processor.Eval(req, mock)

	if mock.Response.Body != expectedResult {
		t.Error("The result differs from the expected result", mock.Response.Body, expectedResult)
	}
}
