package ini

import (
	"fmt"
	"reflect"
	"testing"
)

type TestElement struct {
	inputKey string
	expected []string
}

func TestDecomposeKeyWithDotNotation(t *testing.T) {
	dataProvider := []TestElement{
		{"key", []string{"key"}},
		{"key.foo", []string{"key", "foo"}},
		{"key.foo.bar", []string{"key", "foo", "bar"}},
	}

	for _, element := range dataProvider {
		t.Run("it decomposes "+element.inputKey, func(t *testing.T) {
			actual := DecomposeKeyWithDotNotation(element.inputKey)

			if !reflect.DeepEqual(element.expected, actual) {
				t.Error(fmt.Sprintf("Expected %s, got %s", element.expected, actual))
			}
		})
	}
}

func TestDecomposeKeyWithBracketNotation(t *testing.T) {
	dataProvider := []TestElement{
		{"key", []string{"key"}},
		{"key[foo]", []string{"key", "foo"}},
		{"key[foo][bar]", []string{"key", "foo", "bar"}},
	}

	for _, element := range dataProvider {
		t.Run("it decomposes "+element.inputKey, func(t *testing.T) {
			actual := DecomposeKeyWithBracketNotation(element.inputKey)
			t.Log(actual)

			if !reflect.DeepEqual(element.expected, actual) {
				t.Error(fmt.Sprintf("Expected %s, got %s", element.expected, actual))
			}
		})
	}
}
