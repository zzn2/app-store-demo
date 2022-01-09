package filter

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	var tests = []struct {
		testName string
		input    Rule
		want     string
	}{
		{
			"name=Tom",
			Rule{FieldName: "name", Op: Equals, Value: "Tom"},
			"Rule: name == Tom",
		},
		{
			"name[like]=Tom",
			Rule{FieldName: "name", Op: Like, Value: "Tom"},
			"Rule: name like Tom",
		},
		{
			"age[lt]=20",
			Rule{FieldName: "age", Op: LessThan, Value: "20"},
			"Rule: age < 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			str := fmt.Sprintf("%s", tt.input)
			if str != tt.want {
				t.Errorf("Expect '%s' but got '%s'.", tt.want, str)
			}
		})
	}
}
