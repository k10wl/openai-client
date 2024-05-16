package client

import (
	"reflect"
	"testing"
)

func TestGetModels(t *testing.T) {
	m := NewOpenAIClient("secret api key")
	m.modelsLimits = map[string]limit{
		"foo": {
			name:   "foo",
			tokens: 1,
		},
		"bar": {
			name:   "bar",
			tokens: 1,
		},
	}
	expected := []string{"foo", "bar"}
	actual := m.GetModels()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Failed to get models\nexpected: %v\nactual:   %v",
			expected,
			actual,
		)
	}
}
