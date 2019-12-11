package utils

import (
	"encoding/json"
	"errors"
	"github.com/Jeffail/gabs"
)

//ErrorPropertyMissingInJSON when there's no such property in the JSON document
var ErrorPropertyMissingInJSON = errors.New("There is no such property in the JSON document")

//IsJSON checks if a string is valid JSON or not
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func IsArray(data interface{}) bool {
	_, ok := data.([]interface{})

	return ok
}

func IsObject(data interface{}) bool {
	_, ok := data.(map[string]interface{})

	return ok
}

//JoinJSON merges two JSON strings
func JoinJSON(inputs ...*gabs.Container) *gabs.Container {
	if len(inputs) == 1 {
		return inputs[0]
	}

	result := gabs.New()

	for _, input := range inputs {
		if IsObject(input.Data()) {
			children, _ := input.ChildrenMap()

			for name, child := range children {
				value := child.Data()
				existing := result.S(name)
				existingValue := existing.Data()

				if IsObject(existingValue) && IsObject(value) || IsArray(existingValue) && IsArray(value) {
					result.Set(JoinJSON(existing, child).Data(), name)
				} else {
					result.Set(value, name)
				}
			}

			continue
		}

		if IsArray(input.Data()) {
			array := gabs.New()
			name := "array"
			children, _ := input.Children()
			array.Array(name)

			for k, child := range children {
				value := child.Data()
				existing := result.Index(k)
				existingValue := existing.Data()

				if IsObject(existingValue) && IsObject(value) || IsArray(existingValue) && IsArray(value) {
					array.SetIndex(JoinJSON(existing, child).Data(), k)
				} else {
					array.ArrayAppend(value, name)
				}
			}

			result = array.S(name)
		}
	}

	return result
}
