package vars

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/vtrifonov/http-api-mock/definition"
	"github.com/vtrifonov/http-api-mock/logging"
	"github.com/vtrifonov/http-api-mock/vars/fakedata"
)

var errMissingParameterValue = errors.New("The requested method needs input parameters which are not supplied!")

//FakeVarsFiller parses the data looking for fake data tags or request data tags
type FakeVarsFiller struct {
	Fake fakedata.DataFaker
}

func (fvf FakeVarsFiller) call(data reflect.Value, name string) (string, error) {
	// get a reflect.Value for the method
	methodVal := data.MethodByName(name)
	// turn that into an interface{}
	methodIface := methodVal.Interface()

	typeOfFunction := reflect.TypeOf(methodIface)
	inputParamsCount := typeOfFunction.NumIn()
	// check whether the method has no input parameters
	if inputParamsCount > 0 {
		return "", errMissingParameterValue
	}

	// turn that into a function that has the expected signature
	method := methodIface.(func() string)

	// call the method directly
	res := method()
	return res, nil
}

func (fvf FakeVarsFiller) callWithIntParameter(data reflect.Value, name string, parameter int) string {
	// get a reflect.Value for the method
	methodVal := data.MethodByName(name)
	// turn that into an interface{}
	methodIface := methodVal.Interface()
	// turn that into a function that has the expected signature
	method := methodIface.(func(int) string)
	// call the method directly
	res := method(parameter)
	return res
}

func (fvf FakeVarsFiller) callMethod(name string) (string, bool) {
	method, parameter, hasParameter := fvf.getMethodAndParameter(name)
	if hasParameter {
		name = method
	}

	found := false
	data := reflect.ValueOf(fvf.Fake)
	typ := data.Type()
	if nMethod := data.Type().NumMethod(); nMethod > 0 {
		for i := 0; i < nMethod; i++ {
			method := typ.Method(i)
			if strings.ToLower(method.Name) == strings.ToLower(name) {
				found = true // we found the name regardless
				// does receiver type match? (pointerness might be off)
				if typ == method.Type.In(0) {
					if hasParameter {
						return fvf.callWithIntParameter(data, method.Name, parameter), found
					}

					result, err := fvf.call(data, method.Name)
					if err != nil {
						logging.Printf(err.Error())
					}
					return result, err == nil
				}
			}
		}
	}
	return "", found
}

func (fvf FakeVarsFiller) getMethodAndParameter(input string) (method string, parameter int, success bool) {
	r := regexp.MustCompile(`(?P<method>\w+)\((?P<parameter>.*?)\)`)

	match := r.FindStringSubmatch(input)
	result := make(map[string]string)
	names := r.SubexpNames()
	if len(match) >= len(names) {
		for i, name := range names {
			if i != 0 {
				result[name] = match[i]
			}
		}
	}

	method, success = result["method"]
	if !success {
		return
	}

	parameterString, success := result["parameter"]

	parameter, err := strconv.Atoi(parameterString)
	if err != nil {
		success = false
	}

	return
}

func (fvf FakeVarsFiller) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{\s*fake\.([^{]+?)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if i := strings.Index(tag, "fake."); i == 0 {
			s, found = fvf.callMethod(tag[5:])
		}

		if !found {
			return raw
		}
		return s
	})
}
