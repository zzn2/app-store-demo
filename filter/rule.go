package filter

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

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
	Value     string
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
//    name[like]=Tom   -> name likes "Tom"
//    name=Tome        -> name is exactly "Tom"
//    age[gt]=25       -> age > 25
//
func NewRule(nameAndOp string, value string) (*Rule, error) {
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

// Match checks whether the given obj satisfies the rule.
func (r *Rule) Evaluate(value string) (bool, error) {
	// TODO: How to support generic type for input value?
	// TODO: This part can be moved into Operator.
	switch r.Op {
	case op.Equals:
		return r.Value == value, nil
	case op.Like:
		return strings.Contains(value, string(r.Value)), nil
	default:
		return false, errors.New(fmt.Sprintf("Operator %s currently unsupported.", r.Op))
	}
}

func getStringFieldValue(obj interface{}, fieldName string) (string, error) {
	// TODO: support non-case sensitive case
	if reflect.ValueOf(obj).Kind() == reflect.Struct {
		typ := reflect.TypeOf(obj)
		_, exists := typ.FieldByName(fieldName)
		if !exists {
			return "", errors.New(fmt.Sprintf("Field '%s' does not exist in struct %s.", fieldName, typ))
		}

		v := reflect.ValueOf(obj)
		field := v.FieldByName(fieldName)
		kind := field.Kind()
		if kind != reflect.String {
			return "", errors.New(fmt.Sprintf("Field '%s' is with '%s' type. Currently only supports 'string' fileds.", fieldName, kind))
		}

		return field.String(), nil
	} else {
		return "", errors.New(fmt.Sprintf("Obj '%s' is with '%T' type. Currently only supports structs.", obj, obj))
	}
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule: %s %s %s", r.FieldName, r.Op, r.Value)
}
