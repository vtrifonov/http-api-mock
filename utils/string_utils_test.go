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
