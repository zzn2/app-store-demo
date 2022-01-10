// Package filter provides primitives for filtering objects.
package filter

import (
	"errors"
	"fmt"
	"strings"
)

type RuleSet struct {
	Rules []Rule
}

// Create accepts a map generated from query params and parse them
// as a set of rules and build them inside the RuleSet.
func Create(queryParams map[string][]string) (*RuleSet, error) {
	flt := &RuleSet{}
	for key, value := range queryParams {
		if len(value) > 1 {
			return nil, errors.New(fmt.Sprintf("Key '%s' appeared multiple times with values of '%s'. Currently this case is not unsupported.", key, strings.Join(value, ", ")))
		}

		rule, err := NewRule(key, value[0])
		if err != nil {
			return nil, err
		}
		flt.addRule(*rule)
	}
	return flt, nil
}

// addRule adds a new rule to the given RuleSet.
func (f *RuleSet) addRule(rule Rule) {
	f.Rules = append(f.Rules, rule)
}

// Match checks whether the given obj matches the filter f.
// It returns true if matches, otherwise returns false.
func (f *RuleSet) Match(valueProvider func(string) (string, error)) (bool, error) {
	for _, rule := range f.Rules {
		match, err := rule.Evaluate("")
		if err != nil {
			return false, err
		}

		if !match {
			return false, nil
		}
	}

	return true, nil
}

func (f RuleSet) String() string {
	return fmt.Sprintf("RuleSet: %s", f.Rules)
}
