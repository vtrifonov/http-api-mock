package utils

import "testing"

func TestStringUtils_JoinJSON(t *testing.T) {
	json1 := "{ \"a\": 1, \"b\": 2 }"
	json2 := "{ \"a\": 3, \"c\": 4 }"

	result := JoinJSON(json1, json2)
	expectedResult := "{ \"a\": 3, \"b\": 2, \"c\": 4 }"

	if equal, err := JSONSStringsAreEqual(result, expectedResult); !equal || err != nil {
		t.Error("The result differs from the expected result", result, expectedResult)
	}
}

func TestStringUtils_JoinJSON_ComplexJSON(t *testing.T) {
	json1 := "{ \"a\": { \"aa\": 1, \"ab\": 2}, \"b\": 3 }"
	json2 := "{ \"a\": { \"aa\": 4, \"ac\": 5}, \"c\": 6 }"
	json3 := "{ \"a\": { \"aa\": 7, \"ad\": 8}, \"d\": 9 }"

	result := JoinJSON(json1, json2, json3)
	expectedResult := "{ \"a\": { \"aa\": 7, \"ab\": 2, \"ac\": 5, \"ad\": 8}, \"b\": 3, \"c\": 6, \"d\": 9 }"

	if equal, err := JSONSStringsAreEqual(result, expectedResult); !equal || err != nil {
		t.Error("The result differs from the expected result", result, expectedResult)
	}
}

func TestStringUtils_GetJSONProperty_Simple(t *testing.T) {
	json := "{ \"a\": 1, \"b\": 3 }"

	result, _ := GetJSONProperty(json, "a")

	if result != "1" {
		t.Error("The result differs from the expected result", result, "1")
	}
}

func TestStringUtils_GetJSONProperty_MissingProperty(t *testing.T) {
	json := "{ \"a\": 1, \"b\": 3 }"

	result, err := GetJSONProperty(json, "c")

	if err != ErrorPropertyMissingInJSON {
		t.Error("The result error differs from the expected result", err, ErrorPropertyMissingInJSON)
	}

	if result != "" {
		t.Error("The result differs from the expected result", result, "")
	}
}

func TestStringUtils_GetJSONProperty_Complex(t *testing.T) {
	json := "{ \"a\": { \"aa\": 1, \"ab\": 2}, \"b\": 3 }"

	result, _ := GetJSONProperty(json, "a")
	expectedResult := "{ \"aa\": 1, \"ab\": 2}"

	if equal, err := JSONSStringsAreEqual(result, expectedResult); !equal || err != nil {
		t.Error("The result differs from the expected result", result, expectedResult)
	}
}

func TestStringUtils_GetJSONCompositePropertyValue_Simple(t *testing.T) {
	json := "{ \"a\": 1, \"b\": 3 }"

	result, _ := GetJSONCompositePropertyValue(json, "a")

	if result != "1" {
		t.Error("The result differs from the expected result", result, "1")
	}
}

func TestStringUtils_GetJSONCompositePropertyValue_Complex(t *testing.T) {
	json := "{ \"a\": { \"aa\": 4, \"ab\": 2}, \"b\": 3 }"

	result, _ := GetJSONCompositePropertyValue(json, "a.aa")

	if result != "4" {
		t.Error("The result differs from the expected result", result, "4")
	}
}

func TestStringUtils_GetJSONCompositePropertyValue_Triple(t *testing.T) {
	json := "{ \"a\": { \"aa\": { \"aaa\": 4 }, \"ab\": 2}, \"b\": 3 }"

	result, _ := GetJSONCompositePropertyValue(json, "a.aa.aaa")

	if result != "4" {
		t.Error("The result differs from the expected result", result, "4")
	}
}

func TestStringUtils_GetJSONCompositePropertyValue_Missing(t *testing.T) {
	json := "{ \"a\": { \"aa\": 4, \"ab\": 2}, \"b\": 3 }"

	result, err := GetJSONProperty(json, "a.ac")

	if err != ErrorPropertyMissingInJSON {
		t.Error("The result error differs from the expected result", err, ErrorPropertyMissingInJSON)
	}

	if result != "" {
		t.Error("The result differs from the expected result", result, "")
	}
}

func TestStringUtils_GetPropertyValue_JSON(t *testing.T) {
	input := "{ \"a\": { \"aa\": 4, \"ab\": 2}, \"b\": 3 }"

	result, _ := GetPropertyValue(input, "a.aa")

	if result != "4" {
		t.Error("The result differs from the expected result", result, "4")
	}
}

func TestStringUtils_GetPropertyValue_QueryStrings(t *testing.T) {
	input := "type=smtp&name=My%20New%20Check&resolution=15&sendtoemail=true&sendtosms=true&sendnotificationwhendown=1&contactids=123456,789012&host=smtp.mymailserver.com&auth=myuser%3Amypassword&encryption=true"

	result, _ := GetPropertyValue(input, "name")
	expectedResult := "My New Check"

	if result != expectedResult {
		t.Error("The result differs from the expected result", result, expectedResult)
	}
}
