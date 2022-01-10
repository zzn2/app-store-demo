package op

import (
	"errors"
	"fmt"
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
type Operator string

const (
	Equals      Operator = "=="
	Like                 = "like"
	LessThan             = "<"
	GreaterThan          = ">"
	Unknown              = "Unknown"
)

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
