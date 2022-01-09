package filter

import (
	"fmt"
	"strings"
	"testing"

	"example.com/goserver/filter/op"
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

func TestNew(t *testing.T) {
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
			keyAndValue := strings.Split(tt.input, "=")
			rule, err := New(keyAndValue[0], keyAndValue[1])
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
