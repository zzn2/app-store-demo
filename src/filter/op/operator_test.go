package op

import (
	"fmt"
	"testing"
)

type MyStruct struct {
	Field1 int
	Field2 string
}

type Version struct {
	Major int
	Minor int
}

type ComparableVersion struct {
	Major int
	Minor int
}

func (v ComparableVersion) LessThan(another interface{}) bool {
	if v.Major < another.(ComparableVersion).Major {
		return true
	} else if v.Minor < another.(ComparableVersion).Minor {
		return true
	}

	return false
}

func (v ComparableVersion) GreaterThan(another interface{}) bool {
	if v.Major > another.(ComparableVersion).Major {
		return true
	} else if v.Minor > another.(ComparableVersion).Minor {
		return true
	}

	return false
}

func TestParse(t *testing.T) {
	var tests = []struct {
		input        string
		want         Operator
		errorMessage string
	}{
		{"", Equals, ""},
		{"like", Like, ""},
		{"Like", Like, ""},
		{"LIKE", Like, ""},
		{"lt", LessThan, ""},
		{"gt", GreaterThan, ""},
		{"other", Unknown, "Unrecognized operator type 'other'"},
	}

	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			op, err := Parse(tt.input)
			if op.Name != tt.want.Name {
				t.Errorf("Expect '%s' but got '%s'.", tt.want.Name, op.Name)
			}
			if err != nil {
				if err.Error() != tt.errorMessage {
					t.Errorf("Expect error message '%s' but got '%s'.", tt.errorMessage, err.Error())
				}
			}
		})
	}
}

func TestIsValidType(t *testing.T) {
	var i int
	var s string
	var v Version
	var cv ComparableVersion

	var tests = []struct {
		op       Operator
		value    interface{}
		expected bool
	}{
		{Equals, i, true},
		{Equals, s, true},
		{Equals, v, true},
		{Equals, cv, true},
		{Like, i, false},
		{Like, s, true},
		{Like, v, false},
		{Like, cv, false},
		{LessThan, i, true},
		{LessThan, s, false},
		{LessThan, v, false},
		{LessThan, cv, true},
		{GreaterThan, i, true},
		{GreaterThan, s, false},
		{GreaterThan, v, false},
		{GreaterThan, cv, true},
	}

	for _, tt := range tests {
		t.Run(tt.op.Name, func(t *testing.T) {
			valid := tt.op.IsValidType(tt.value)
			if tt.expected != valid {
				t.Errorf("Expected to be '%v' but got '%v'", tt.expected, valid)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	var tests = []struct {
		operator      Operator
		incomingValue interface{}
		threshold     interface{}
		want          bool
		errorMessage  string
	}{
		{Equals, 0, 0, true, ""},
		{Equals, 0, 1, false, ""},
		{Equals, "a", "a", true, ""},
		{Equals, "a", "b", false, ""},
		{Equals, MyStruct{Field1: 42, Field2: "Hello"}, MyStruct{Field1: 42, Field2: "Hello"}, true, ""},
		{Equals, MyStruct{Field1: 42, Field2: "Hello"}, MyStruct{Field1: 0, Field2: "Hello"}, false, ""},
		{Equals, MyStruct{Field1: 42, Field2: "Hello"}, MyStruct{Field1: 42, Field2: "World"}, false, ""},
		{Like, "abcde", "abc", true, ""},
		{Like, "abc", "abc", true, ""},
		{Like, "abc", "abcde", false, ""},
		{Like, 1, 2, false, "Operator 'Like' does not support the incoming values in int type."},
		{LessThan, 1, 2, true, ""},
		{LessThan, 2, 1, false, ""},
		{LessThan, Version{Major: 1, Minor: 0}, Version{Major: 1, Minor: 1}, false, "Operator 'LessThan' does not support the incoming values in op.Version type."},
		{LessThan, ComparableVersion{Major: 1, Minor: 0}, ComparableVersion{Major: 1, Minor: 1}, true, ""},
		{GreaterThan, 1, 2, false, ""},
		{GreaterThan, 2, 1, true, ""},
		{GreaterThan, Version{Major: 1, Minor: 0}, Version{Major: 1, Minor: 1}, false, "Operator 'GreaterThan' does not support the incoming values in op.Version type."},
		{GreaterThan, ComparableVersion{Major: 1, Minor: 0}, ComparableVersion{Major: 1, Minor: 1}, false, ""},
		{GreaterThan, ComparableVersion{Major: 1, Minor: 0}, 1, false, "TypeMismatch: Expects incoming value to be 'int' type but was 'op.ComparableVersion'"},
	}

	for _, tt := range tests {
		testname := tt.operator.Name
		t.Run(testname, func(t *testing.T) {
			res, err := tt.operator.Evaluate(tt.incomingValue, tt.threshold)
			if res != tt.want {
				t.Errorf("Expect to be '%v' but got '%v'", tt.want, res)
			}
			if err != nil {
				if err.Error() != tt.errorMessage {
					t.Errorf("Expect error message '%s' but got '%s'.", tt.errorMessage, err.Error())
				}
			}
		})
	}
}

func TestIsNumberType(t *testing.T) {
	var v Version
	var cv ComparableVersion

	var tests = []struct {
		value interface{}
		want  bool
	}{
		{1, true},
		{"", false},
		{1.1, true},
		{true, false},
		{v, false},
		{cv, false},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.value)
		t.Run(testName, func(t *testing.T) {
			actual := isNumberType(tt.value)
			if tt.want != actual {
				t.Errorf("isNumberType(%v) expected to be %v but got %v", tt.value, tt.want, actual)
			}
		})
	}
}

func TestIsStringType(t *testing.T) {
	var v Version
	var cv ComparableVersion

	var tests = []struct {
		value interface{}
		want  bool
	}{
		{1, false},
		{"", true},
		{1.1, false},
		{true, false},
		{v, false},
		{cv, false},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.value)
		t.Run(testName, func(t *testing.T) {
			actual := isStringType(tt.value)
			if tt.want != actual {
				t.Errorf("isStringType(%v) expected to be %v but got %v", tt.value, tt.want, actual)
			}
		})
	}
}

func TestIsValueComparerType(t *testing.T) {
	var v Version
	var cv ComparableVersion

	var tests = []struct {
		value interface{}
		want  bool
	}{
		{1, false},
		{"", false},
		{1.1, false},
		{true, false},
		{v, false},
		{cv, true},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.value)
		t.Run(testName, func(t *testing.T) {
			actual := isValueComparerType(tt.value)
			if tt.want != actual {
				t.Errorf("isValueComparerType(%v) expected to be %v but got %v", tt.value, tt.want, actual)
			}
		})
	}
}
