package urlbuilder

import (
	"reflect"
	"testing"
)

func TestURLBuilder_BuildString(t *testing.T) {
	// Set up test data
	u := URLBuilder{
		root_uri:  "https://example.com",
		arguments: map[string]string{"key1": "value1", "key2": "value2"},
	}

	// Call the function being tested
	result := u.BuildString()

	// Assert the expected value
	expected := "https://example.com?key1=value1&key2=value2"
	if result != expected {
		t.Errorf("Expected URL to be '%s', but got '%s'", expected, result)
	}
}

func TestURLBuilder_AddArgument(t *testing.T) {
	// Set up test data
	u := URLBuilder{
		root_uri:  "https://example.com",
		arguments: make(map[string]string),
	}

	// Call the function being tested
	u.AddArgument("key1", "value1").
		AddArgument("key2", "value2")

	// Assert the expected values
	expected := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(u.arguments, expected) {
		t.Errorf("Expected arguments to be '%v', but got '%v'", expected, u.arguments)
	}
}
