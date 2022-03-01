package filter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/zzn2/demo/appstore/filter/op"
	"github.com/zzn2/demo/appstore/semver"
)

type User struct {
	FirstName string
	LastName  string
	Age       int
	Version   semver.Version
}

func TestParseRule(t *testing.T) {
	var u User
	var tests = []struct {
		input          string
		expectedRule   Rule
		expectedErrMsg string
	}{
		{
			"FirstName=Tom",
			Rule{FieldName: "FirstName", Op: op.Equals, Value: "Tom"},
			"",
		},
		{
			"firstname=Tom",
			Rule{FieldName: "firstname", Op: op.Equals, Value: "Tom"},
			"",
		},
		{
			"name=Tom",
			Rule{},
			"Failed to create rule: Field with name 'name' does not exist.",
		},
		{
			"FirstName[like]=Tom",
			Rule{FieldName: "FirstName", Op: op.Like, Value: "Tom"},
			"",
		},
		{
			"age[lt]=20",
			Rule{FieldName: "age", Op: op.LessThan, Value: 20},
			"",
		},
		{
			"age[dummy]=20",
			Rule{},
			"Unrecognized operator type 'dummy'",
		},
		{
			"version[lt]=0.0.1",
			Rule{FieldName: "version", Op: op.LessThan, Value: semver.Version{Major: 0, Minor: 0, Patch: 1}},
			"",
		},
		{
			"version[lt]=0.0.a",
			Rule{},
			`Failed to create rule: Failed to parse version '0.0.a': Invalid character(s) found in number "a"`,
		},
		{
			"version[like]=0.0.1",
			Rule{},
			"Failed to create rule: Type 'semver.Version' does not support 'Like' operator",
		},
		{
			"illegal/keyformat=1",
			Rule{},
			"Malformed input key format: 'illegal/keyformat'",
		},
	}

	for _, tt := range tests {
		testName := tt.input
		t.Run(testName, func(t *testing.T) {
			rule, err := ParseRule(tt.input, u)
			if rule != tt.expectedRule {
				t.Errorf("Expect '%s' but got '%s'.", tt.expectedRule, rule)
			}
			if err != nil {
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expect err to be '%s' but got '%s'.", tt.expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	var u User
	var tests = []struct {
		ruleText        string
		valueToEvaluate interface{}
		expectedResult  bool
		expectedErrMsg  string
	}{
		{
			"FirstName=Tom",
			"Tom",
			true,
			"",
		},
		{
			"LastName=Tom",
			"Green",
			false,
			"",
		},
		{
			"LastName=ree",
			"Green",
			false,
			"",
		},
		{
			"LastName[Like]=ree",
			"Green",
			true,
			"",
		},
		{
			"LastName[Like]=Gre",
			"Green",
			true,
			"",
		},
		{
			"LastName[Like]=een",
			"Green",
			true,
			"",
		},
		{
			"LastName[Like]=Green",
			"Green",
			true,
			"",
		},
		{
			"Age[lt]=25",
			20,
			true,
			"",
		},
		{
			"Age[gt]=25",
			20,
			false,
			"",
		},
		{
			"Age[gt]=25",
			"20",
			false,
			"Failed to evaluate 'Rule: Age gt 25 (int)': TypeMismatch: Expects incoming value to be 'int' type but was 'string'",
		},
	}

	for _, tt := range tests {
		testName := strings.NewReplacer("[", "_", "]", "_").Replace(tt.ruleText)
		t.Run(testName, func(t *testing.T) {
			rule, _ := ParseRule(tt.ruleText, u)
			match, err := rule.Evaluate(tt.valueToEvaluate)
			if match != tt.expectedResult {
				t.Errorf("Expect '%v' but got '%v'.", tt.expectedResult, match)
			}
			if err != nil {
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expect err to be '%s' but got '%s'.", tt.expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestGetNameAndOp(t *testing.T) {
	var tests = []struct {
		testName       string
		input          string
		expectedName   string
		expectedOp     op.Operator
		expectedErrMsg string
	}{
		{
			"Normal param name without brackets",
			"param",
			"param",
			op.Equals,
			"",
		},
		{
			"Param name with empty LHS bracket",
			"param[]",
			"param",
			op.Equals,
			"",
		},
		{
			"Param name with LHS brackets (op.Like)",
			"param[like]",
			"param",
			op.Like,
			"",
		},
		{
			"Param name with LHS brackets (op.LessThan)",
			"param[lt]",
			"param",
			op.LessThan,
			"",
		},
		{
			"Param name with LHS brackets (invalid op)",
			"param[dummy]",
			"param",
			op.Unknown,
			"Unrecognized operator type 'dummy'",
		},
		{
			"Illegal input text",
			"param[",
			"",
			op.Unknown,
			"Malformed input key format: 'param['",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			name, op, err := getNameAndOp(tt.input)
			if name != tt.expectedName {
				t.Errorf("Expect name to be '%s' but got '%s'.", tt.expectedName, name)
			}
			if op != tt.expectedOp {
				t.Errorf("Expect operator to be '%s' but got '%s'.", tt.expectedOp, op)
			}
			if err != nil {
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expect err to be '%s' but got '%s'.", tt.expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestGetFieldByName(t *testing.T) {
	var u User
	var tests = []struct {
		testName          string
		obj               interface{}
		fieldName         string
		expectedToBeValid bool
	}{
		{
			"Get field name",
			u,
			"FirstName",
			true,
		},
		{
			"Get field name",
			u,
			"firstname",
			true,
		},
		{
			"Get field name",
			u,
			"Age",
			true,
		},
		{
			"Get field name",
			u,
			"Unknown",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			res := getFieldByName(tt.obj, tt.fieldName)
			valid := res.IsValid()
			if valid != tt.expectedToBeValid {
				t.Errorf("IsValid() expected to be '%v' but got '%v'", tt.expectedToBeValid, valid)
			}
		})
	}
}

func TestParseText(t *testing.T) {
	var i int
	var i32 int32
	var i64 int64
	var u uint
	var u32 uint32
	var u64 uint64
	var s string
	var v semver.Version

	var tests = []struct {
		testName             string
		textToParse          string
		parseAsType          reflect.Type
		expectedResult       interface{}
		expectedErrorMessage string
	}{
		{
			"int",
			"42",
			reflect.TypeOf(i),
			int(42),
			"",
		},
		{
			"int32",
			"42",
			reflect.TypeOf(i32),
			int32(42),
			"",
		},
		{
			"int64",
			"42",
			reflect.TypeOf(i64),
			int64(42),
			"",
		},
		{
			"int64_error",
			"a",
			reflect.TypeOf(i64),
			int64(0),
			"Invalid integer format: strconv.ParseInt: parsing \"a\": invalid syntax",
		},
		{
			"uint",
			"42",
			reflect.TypeOf(u),
			uint(42),
			"",
		},
		{
			"uint32",
			"42",
			reflect.TypeOf(u32),
			uint32(42),
			"",
		},
		{
			"uint64",
			"42",
			reflect.TypeOf(u64),
			uint64(42),
			"",
		},
		{
			"uint64_error",
			"-1",
			reflect.TypeOf(u64),
			uint64(0),
			"Invalid integer format: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			"string",
			"hello world",
			reflect.TypeOf(s),
			"hello world",
			"",
		},
		{
			"version",
			"0.0.1",
			reflect.TypeOf(v),
			semver.Version{Major: 0, Minor: 0, Patch: 1},
			"",
		},
		{
			"version_error",
			"0.0.a",
			reflect.TypeOf(v),
			semver.Version{Major: 0, Minor: 0, Patch: 0},
			"Failed to parse version '0.0.a': Invalid character(s) found in number \"a\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			parsed, err := parseText(tt.textToParse, tt.parseAsType)
			if err != nil {
				if err.Error() != tt.expectedErrorMessage {
					t.Errorf("Expected error message '%s' but got '%s'", tt.expectedErrorMessage, err.Error())
				}
			}
			if parsed != tt.expectedResult {
				t.Errorf("Expected '%v' (%T) but got '%v' (%T)", tt.expectedResult, tt.expectedResult, parsed, parsed)
			}
		})
	}
}
