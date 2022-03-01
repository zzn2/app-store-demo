// Package op provides operators used in filter evaluation.
package op

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Operator defines the operation type in a filter rule.
// Typically it is defined in a query string using LHS brackets like:
//
//    /employees?name[like]=Tom&age[gt]=25
//
// When there's no LHS brackets, it represents a precise match:
//
// .  /employees?name=Tommy
//
type Operator struct {
	Name   string
	Symbol string
	OpText string
}

var (
	Equals = Operator{
		Name:   "Equals",
		Symbol: "==",
		OpText: "eq",
	}

	Like = Operator{
		Name:   "Like",
		Symbol: "like",
		OpText: "like",
	}

	LessThan = Operator{
		Name:   "LessThan",
		Symbol: "<",
		OpText: "lt",
	}

	GreaterThan = Operator{
		Name:   "GreaterThan",
		Symbol: ">",
		OpText: "gt",
	}

	// More operators can be added here.

	Unknown = Operator{}
)

// LessThanComparer describes the operation to compare whether less than a given value.
type LessThanComparer interface {
	LessThan(another interface{}) bool
}

// GreaterThanComparer describes the operation to compare whether greater than a given value.
type GreaterThanComparer interface {
	GreaterThan(another interface{}) bool
}

// ValueComparer describes operations of comparing two given objects.
type ValueComparer interface {
	LessThanComparer
	GreaterThanComparer
}

// Parse parses given text into Operation.
// It returns the corresponding Operation object when parse succeed
// and generates error when the operation is unrecognized.
func Parse(text string) (Operator, error) {
	switch strings.ToLower(text) {
	case "":
		return Equals, nil
	case "like":
		return Like, nil
	case "lt":
		return LessThan, nil
	case "gt":
		return GreaterThan, nil
	default:
		return Unknown, errors.New(fmt.Sprintf("Unrecognized operator type '%s'", text))
	}
}

// IsValidType returns whether a given object is valid for this operator.
// e.g.
//   1. When the operator is "Like", it only accepts value in string type.
//   2. When the operator is "lt" or "gt", it accepts numbers or comparable objects, but no strings.
//   3. etc.
func (op Operator) IsValidType(value interface{}) bool {
	switch op {
	case Equals:
		// 'Equals' operator is available for any type
		return true
	case Like:
		// 'Like' operator only valid for string types
		return isStringType(value)
	case LessThan, GreaterThan:
		// 'LessThan', 'GreaterThan' operators can accept either a number
		// or an object which implements 'ValueComparer' interface.
		return isNumberType(value) || isValueComparerType(value)
	}

	return false
}

// Evaluate applies the incomingValue to the operator and baseValue.
// e.g.
//   For LessThan operator, given incomingValue=5, baseValue=10, will return true since 5<10.
//   The baseValue is typically set by a filter.Rule object,
//   which is parsed from querystring like: age[lt]=10
//   And the incomingValue is the value of the user's input, in this example the value is 5.
func (op Operator) Evaluate(incomingValue interface{}, baseValue interface{}) (bool, error) {
	// Make sure incomingValue is the same type with baseValue.
	if reflect.ValueOf(incomingValue).Type() != reflect.ValueOf(baseValue).Type() {
		return false, fmt.Errorf("TypeMismatch: Expects incoming value to be '%T' type but was '%T'", baseValue, incomingValue)
	}
	// Usually we have already verified the type of baseValue to be compatible to the operator.
	// And we have already made sure the types of incomingValue and baseValue are the same.
	// So, the type of incomingValue must be compatible with the operator.
	// However, we do another type check here to make sure for the assumption.
	if !op.IsValidType(incomingValue) {
		return false, fmt.Errorf("Operator '%s' does not support the incoming values in %T type.", op, incomingValue)
	}

	switch op {
	case Equals:
		return incomingValue == baseValue, nil
	case Like:
		return strings.Contains(incomingValue.(string), baseValue.(string)), nil
	case LessThan:
		if isNumberType(baseValue) {
			// Assume all number types are 'int'.
			// TODO: adapt to all the number types.
			return incomingValue.(int) < baseValue.(int), nil
		} else if isValueComparerType(baseValue) {
			return incomingValue.(ValueComparer).LessThan(baseValue.(ValueComparer)), nil
		} else {
			panic("Should never fall into this branch.")
		}
	case GreaterThan:
		if isNumberType(baseValue) {
			return incomingValue.(int) > baseValue.(int), nil
		} else if isValueComparerType(baseValue) {
			return incomingValue.(ValueComparer).GreaterThan(baseValue.(ValueComparer)), nil
		} else {
			panic("Should never fall into this branch.")
		}
	default:
		return false, fmt.Errorf("Unknown operator '%s'", op)
	}
}

// String returns a text representation of this object.
func (op Operator) String() string {
	return op.Name
}

func isNumberType(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	}

	return false
}

func isStringType(value interface{}) bool {
	switch value.(type) {
	case string:
		return true
	}

	return false
}

func isValueComparerType(value interface{}) bool {
	switch value.(type) {
	case ValueComparer:
		return true
	}

	return false
}
