package filter

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/zzn2/demo/appstore/filter/op"
)

// Rule describes a matching rule.
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
	regexForParamWithLhsBracket = regexp.MustCompile(`^(?P<name>[a-zA-Z0-9.]+)\[(?P<op>[a-zA-Z]*)\]$`)
)

// NewRule creates a new instance of Rule.
// The nameAndOp shows the parameter name to be validated and with a operational operator in LHS bracket.
// value refers the value attached to the check rule.
// e.g.
//
//    name[like]=Tom   -> name likes "Tom"
//    name=Tome        -> name is exactly "Tom"
//    age[gt]=25       -> age > 25
//
func NewRule(nameAndOp string, value string, applyToObj interface{}) (Rule, error) {
	name, operator, error := getNameAndOp(nameAndOp)
	if error != nil {
		return Rule{}, error
	}

	field := getFieldByName(applyToObj, name)
	if !field.IsValid() {
		return Rule{}, fmt.Errorf("Failed to create rule: Field with name '%s' does not exist.", name)
	}

	parsedValue, err := parseText(value, field.Type())
	if err != nil {
		return Rule{}, fmt.Errorf("Failed to create rule: %w", err)
	}

	if !operator.IsValidType(parsedValue) {
		return Rule{}, fmt.Errorf("Failed to create rule: Type '%T' does not support '%s' operator", parsedValue, operator)
	}

	return Rule{
		FieldName: name,
		Op:        operator,
		Value:     parsedValue,
	}, nil
}

// ParseRule parses text representation (samples listed below) into new instance of Rule.
//
//    name[like]=Tom   -> name likes "Tom"
//    name=Tome        -> name is exactly "Tom"
//    age[gt]=25       -> age > 25
//
// If failed to parse, it returns nil and the detail error.
func ParseRule(text string, applyToObj interface{}) (Rule, error) {
	keyAndValue := strings.SplitN(text, "=", 2)
	return NewRule(keyAndValue[0], keyAndValue[1], applyToObj)
}

// Match checks whether the given obj satisfies the rule.
func (r Rule) Match(obj interface{}) (bool, error) {
	field := r.getField(obj).Interface()
	return r.Evaluate(field)
}

// Evaluate evaluates whether a given value satisfies the rule.
func (r Rule) Evaluate(value interface{}) (bool, error) {
	succeed, err := r.Op.Evaluate(value, r.Value)
	if err != nil {
		return false, fmt.Errorf("Failed to evaluate '%s': %w", r, err)
	}
	return succeed, err
}

// String returns a text representation for this rule.
func (r Rule) String() string {
	return fmt.Sprintf("Rule: %s %v %v (%T)", r.FieldName, r.Op.OpText, r.Value, r.Value)
}

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

// getField gets the field from the given object with the name of this rule.
func (r Rule) getField(v interface{}) reflect.Value {
	return getFieldByName(v, r.FieldName)
}

// getFieldByName gets the field from given object with specific field name.
func getFieldByName(v interface{}, name string) reflect.Value {
	match := func(fieldName string) bool {
		return strings.EqualFold(name, fieldName)
	}
	return reflect.ValueOf(v).FieldByNameFunc(match)
}

// parseText parses text to object of the given type.
func parseText(text string, asType reflect.Type) (interface{}, error) {
	kind := asType.Kind()
	switch kind {
	case reflect.String:
		return text, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return reflect.Zero(asType).Interface(), fmt.Errorf("Invalid integer format: %w", err)
		}
		result := reflect.New(asType).Elem()
		result.SetInt(parsed)
		return result.Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return reflect.Zero(asType).Interface(), fmt.Errorf("Invalid integer format: %w", err)
		}
		result := reflect.New(asType).Elem()
		result.SetUint(parsed)
		return result.Interface(), nil
	default:
		unmarshaler := reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
		ptrToType := reflect.PtrTo(asType)
		if ptrToType.Implements(unmarshaler) {
			instance := reflect.New(asType).Interface()
			err := instance.(encoding.TextUnmarshaler).UnmarshalText([]byte(text))
			return reflect.ValueOf(instance).Elem().Interface(), err
		}
		return nil, fmt.Errorf("Unable to parse '%s' into given type '%s'", text, asType.Name())
	}
}
