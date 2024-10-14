package ini

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDecomposeKeyWithDotNotation(t *testing.T) {
	type TestElement struct {
		inputKey string
		expected []string
	}

	dataProvider := []TestElement{
		{"key", []string{"key"}},
		{"key.foo", []string{"key", "foo"}},
		{"key.foo.bar", []string{"key", "foo", "bar"}},
	}

	for _, element := range dataProvider {
		t.Run("it decomposes "+element.inputKey, func(t *testing.T) {
		})

		actual := DecomposeKeyWithDotNotation(element.inputKey)

		if !reflect.DeepEqual(element.expected, actual) {
			t.Error(fmt.Sprintf("Expected %s, got %s", element.expected, actual))
		}
	}
}
