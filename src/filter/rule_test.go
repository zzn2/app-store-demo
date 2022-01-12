package filter

import (
	"fmt"
	"testing"

	"github.com/zzn2/demo/appstore/filter/op"
)

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

func TestParseRule(t *testing.T) {
	var tests = []struct {
		input           string
		expectedRuleStr string
		expectedErrMsg  string
	}{
		{
			"name=Tom",
			"Rule: name == Tom",
			"",
		},
		{
			"name[like]=Tom",
			"Rule: name like Tom",
			"",
		},
		{
			"age[lt]=20",
			"Rule: age < 20",
			"",
		},
		{
			"age[dummy]=20",
			"<nil>",
			"Unrecognized operator type 'dummy'",
		},
		{
			"illegal/keyformat=1",
			"<nil>",
			"Malformed input key format: 'illegal/keyformat'",
		},
	}

	for _, tt := range tests {
		testName := tt.input
		t.Run(testName, func(t *testing.T) {
			rule, err := ParseRule(tt.input)
			str := fmt.Sprintf("%s", rule)
			if str != tt.expectedRuleStr {
				t.Errorf("Expect '%s' but got '%s'.", tt.expectedRuleStr, str)
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
	var tests = []struct {
		ruleText        string
		valueToEvaluate string
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
			"20",
			false,
			"Operator '<' currently unsupported.",
		},
		{
			"Age[gt]=25",
			"20",
			false,
			"Operator '>' currently unsupported.",
		},
	}

	for _, tt := range tests {
		testName := tt.ruleText
		t.Run(testName, func(t *testing.T) {
			rule, _ := ParseRule(tt.ruleText)
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

func TestGetStringFieldValue(t *testing.T) {
	type Employee struct {
		FirstName string
		LastName  string
		Age       int
	}
	var tomGreen = Employee{
		FirstName: "Tom",
		LastName:  "Green",
		Age:       25,
	}

	var tests = []struct {
		testName       string
		inputObj       Employee
		fieldName      string
		expectedValue  string
		expectedErrMsg string
	}{
		{
			"Normal case",
			tomGreen,
			"FirstName",
			"Tom",
			"",
		},
		{
			"Normal case 2",
			tomGreen,
			"LastName",
			"Green",
			"",
		},
		{
			"Refer to non-string field",
			tomGreen,
			"Age",
			"",
			"Field 'Age' is with 'int' type. Currently only supports 'string' fileds.",
		},
		{
			"Refer to non-exist field",
			tomGreen,
			"Address",
			"",
			"Field 'Address' does not exist in struct filter.Employee.",
		},
		{
			"Refer to non-struct type",
			tomGreen, // This field is not referenced. It will be replaced to a plain string later.
			"NonStruct",
			"",
			"Obj 'plain string' is with 'string' type. Currently only supports structs.",
		},
		/*
			Skip for now since it is not on the critical path.
			TODO: Enable when possible.
				{
					"Reference with pointer",
					tomGreen, // This field is not referenced. It will be replaced to a plain string later.
					"FirstName",
					"Tom",
					"",
				},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			var obj interface{} = tt.inputObj

			// The value in `tests` slice must be with `Employee` type.
			// However, we do need some test case for non-Employee types.
			// So we are overriding them here for hack.
			// TODO: Find a better way to do so.
			if tt.testName == "Refer to non-struct type" {
				obj = "plain string"
			} else if tt.testName == "Reference with pointer" {
				obj = &tomGreen
			}

			value, err := getStringFieldValue(obj, tt.fieldName)
			if value != tt.expectedValue {
				t.Errorf("Expect '%s' but got '%s'.", tt.expectedValue, value)
			}
			if err != nil {
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expect err to be '%s' but got '%s'.", tt.expectedErrMsg, err.Error())
				}
			}
		})
	}
}
