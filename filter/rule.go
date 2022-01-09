package filter

import (
	"errors"
	"fmt"
	"regexp"

	"example.com/goserver/filter/op"
)

// Rule describes a filter rule.
// Typically it is defined in a query string using LHS brackets like:
//
//    name[like]=Tom -> To find name like "Tom"
//    name=Tome      -> To find name exactly is "Tom"
//    age[gt]=25     -> To find age > 25
//
type Rule struct {
	FieldName string
	Op        op.Operator
	Value     interface{}
}

var (
	regexForPlainParam          = regexp.MustCompile(`^[a-zA-Z0-9.]+$`)
	regexForParamWithLhsBracket = regexp.MustCompile(`^(?P<name>[a-zA-Z0-9]+)\[(?P<op>[a-zA-Z]*)\]$`)
)

// getNameAndOp parses a given text and separate them into name and operator.
// The text is expected to be in following format:
//
//     param        returns param, op.Equals
//     param[like]  returns param, op.Like
//     param[gt]    returns param, op.GreaterThan
//
// For unrecognized operator names or illegal input formats, return error.
func getNameAndOp(text string) (name string, operator op.Operator, err error) {
	bytes := []byte(text)
	if regexForPlainParam.Match(bytes) {
		return text, op.Equals, nil
	} else if regexForParamWithLhsBracket.Match(bytes) {
		// Given "param[like]", returns:
		//   [0]: param[like]
		//   [1]: param
		//   [2]: like
		// So, index 1 will be the param and index 2 will be the operator
		match := regexForParamWithLhsBracket.FindStringSubmatch(text)
		if operator, err = op.Parse(match[2]); err != nil {
			return match[1], operator, err
		} else {
			return match[1], operator, nil
		}
	} else {
		return "", op.Unknown, errors.New(fmt.Sprintf("Malformed input key format: '%s'", text))
	}
}

// New creates a new instance of Rule.
// The nameAndOp shows the parameter name to be validated and with a operational operator in LHS bracket.
// value refers the value attached to the check rule.
// e.g.
//
//    name[like], Tom  -> name likes "Tom"
//    name=Tome        -> name is exactly "Tom"
//    age[gt]=25       -> age > 25
//
func New(nameAndOp string, value string) (*Rule, error) {
	name, operator, error := getNameAndOp(nameAndOp)
	if error != nil {
		return nil, error
	}

	// TODO: Assert the value should represent a number when the op is LessThan or GreaterThan.
	//       I'm not going to implement the details for this demo project since it's not used in the cases.

	return &Rule{
		FieldName: name,
		Op:        operator,
		Value:     value,
	}, nil
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule: %s %s %s", r.FieldName, r.Op, r.Value)
}
