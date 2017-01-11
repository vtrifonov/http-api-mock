package utils

import (
	"encoding/json"
	"errors"
	"net/url"

	"strings"

	"github.com/Jeffail/gabs"
	"github.com/vtrifonov/http-api-mock/definition"
)

//ErrorPropertyMissingInJSON when there's no such property in the JSON document
var ErrorPropertyMissingInJSON = errors.New("There is no such property in the JSON document")

//IsJSON checks if a string is valid JSON or not
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

//JoinJSON merges two JSON strings
func JoinJSON(inputs ...string) string {
	if len(inputs) == 1 {
		return inputs[0]
	}

	result := gabs.New()
	for _, input := range inputs {
		jsonParsed, _ := gabs.ParseJSON([]byte(input))
		children, _ := jsonParsed.S().ChildrenMap()

		for key, child := range children {
			value := child.Data()
			stringValue := child.String()
			// check if the value is again json or a primitive type
			if strings.Index(stringValue, "{") == 0 || strings.Index(stringValue, "[") == 0 {
				// if the value is JSON and we already have value for that key in the result
				existingValue := result.S(key).String()
				if existingValue == "{}" {
					result.Set(value, key)
				} else {
					// merge the inner JSONs
					joinedResult := JoinJSON(existingValue, stringValue)
					valueToStore, _ := gabs.ParseJSON([]byte(joinedResult))
					result.Set(valueToStore.Data(), key)
				}
			} else {
				result.Set(value, key)
			}
		}
	}

	return result.String()
}

//JoinContent returns two contents joined as JSON if both are JSONs otherwise concatenates them
func JoinContent(value1 string, value2 string) string {
	if value1 == "" {
		return value2
	} else if value2 == "" {
		return value1
	} else if (IsJSON(value1)) && IsJSON(value2) {
		return JoinJSON(value1, value2)
	} else {
		return value1 + "\n" + value2
	}
}

//FormatJSON formats a JSON string
func FormatJSON(input string) (result string, err error) {
	var jsonParsed interface{}
	json.Unmarshal([]byte(input), &jsonParsed)
	if err != nil {
		return "", err
	}

	byteString, err := json.Marshal(jsonParsed)
	if err != nil {
		return "", err
	}

	return string(byteString), nil
}

//JSONSStringsAreEqual checks whether two JSON strings are actually equal JSON objects
func JSONSStringsAreEqual(input1 string, input2 string) (result bool, err error) {
	formatedInput1, err := FormatJSON(input1)
	if err != nil {
		return false, err
	}
	formatedInput2, err := FormatJSON(input2)
	if err != nil {
		return false, err
	}
	return formatedInput1 == formatedInput2, nil
}

//WrapNonJSONStringIfNeeded wrapps non JSON string in NonJSONItem object
func WrapNonJSONStringIfNeeded(input string) (result string, err error) {
	if !IsJSON(input) {
		wrapper := definition.NonJSONItem{Content: input}
		bytesString, err := json.Marshal(wrapper)
		if err != nil {
			return "", err
		}
		return string(bytesString), nil
	}
	return input, nil
}

//UnWrapNonJSONStringIfNeeded wrapps non JSON string in NonJSONItem object
func UnWrapNonJSONStringIfNeeded(input string) string {
	if IsJSON(input) && strings.Index(input, "non_json_content") > -1 {
		var nonJSON definition.NonJSONItem
		err := json.Unmarshal([]byte(input), &nonJSON)
		if err != nil || nonJSON.Content == "" { // the json most probably is not a NonJSONItem
			return input
		}

		return nonJSON.Content
	}
	return input
}

//JSONSerialize serializes an inteface to JSON string
func JSONSerialize(input interface{}) (string, error) {
	byteResult, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(byteResult), nil
}

//JSONDeserialize deserializes a JSON string to interface map
func JSONDeserialize(input string) (map[string]interface{}, error) {
	var properties map[string]interface{}
	if err := json.Unmarshal([]byte(input), &properties); err != nil {
		return nil, err
	}
	return properties, nil
}

//GetJSONProperty returns the string value of a given property inside a JSON document
func GetJSONProperty(input string, property string) (string, error) {
	properties, err := JSONDeserialize(input)

	if err != nil {
		return "", err
	}

	result, exists := properties[property]

	if !exists {
		return "", ErrorPropertyMissingInJSON
	}

	return JSONSerialize(result)
}

//GetJSONCompositePropertyValue returns composite property value digging inside complex JSON documents
func GetJSONCompositePropertyValue(input string, property string) (string, error) {
	properties := strings.Split(property, ".")
	value, err := GetJSONProperty(input, properties[0])
	if err != nil {
		return "", err
	}
	if len(properties) > 1 {
		subProperty := strings.Join(properties[1:], ".")
		return GetJSONCompositePropertyValue(value, subProperty)
	}

	value = trimSurroundings(value, "\"")
	value = trimSurroundings(value, "'")

	return value, nil
}

func trimSurroundings(input string, surrounding string) string {
	if strings.HasPrefix(input, surrounding) && strings.HasSuffix(input, surrounding) {
		input = strings.TrimPrefix(input, surrounding)
		input = strings.TrimSuffix(input, surrounding)
	}
	return input
}

//GetPropertyValue returns the json property value if input is json, otherwise tries to parse the value as query string and get property value
func GetPropertyValue(input string, property string) (string, error) {
	if IsJSON(input) {
		return GetJSONCompositePropertyValue(input, property)
	}

	values, err := url.ParseQuery(input)
	if err != nil {
		return "", err
	}
	return values.Get(property), err
}
