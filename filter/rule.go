package filter

import "fmt"

type Rule struct {
	FieldName string
	Op        Operator
	Value     interface{}
}

func (r Rule) String() string {
	return fmt.Sprintf("Rule: %s %s %s", r.FieldName, r.Op, r.Value)
}
