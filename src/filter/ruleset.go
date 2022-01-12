// Package filter provides definitions used for RuleSet to filter objects.
package filter

import (
	"errors"
	"fmt"
	"strings"
)

// RuleSet consists a set of Rules.
type RuleSet struct {
	Rules []Rule
}

// Create accepts a map generated from query params and parse them
// as a set of rules and build them inside the RuleSet.
func CreateRuleSet(queryParams map[string][]string) (*RuleSet, error) {
	flt := &RuleSet{}
	for key, value := range queryParams {
		if len(value) > 1 {
			return nil, errors.New(fmt.Sprintf("Key '%s' appeared multiple times with values of '%s'. Currently this case is not unsupported.", key, strings.Join(value, ", ")))
		}

		rule, err := NewRule(key, value[0])
		if err != nil {
			return nil, err
		}
		flt.AddRule(*rule)
	}
	return flt, nil
}

// AddRule adds a new rule to the given RuleSet.
func (f *RuleSet) AddRule(rule Rule) {
	f.Rules = append(f.Rules, rule)
}

func (f RuleSet) String() string {
	return fmt.Sprintf("RuleSet: %s", f.Rules)
}
