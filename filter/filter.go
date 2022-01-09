// Package filter provides primitives for filtering objects.
package filter

import (
	"errors"
	"fmt"
	"strings"
)

type Filter struct {
	rules []Rule
}

// Create accepts a map generated from query params and parse them
// as a set of rules and build them inside the Filter.
func Create(queryParams map[string][]string) (*Filter, error) {
	flt := &Filter{}
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

// addRule adds a new rule to the given Filter.
func (f *Filter) addRule(rule Rule) {
	f.rules = append(f.rules, rule)
}

// Match checks whether the given obj matches the filter f.
// It returns true if matches, otherwise returns false.
func (f *Filter) Match(obj interface{}) (bool, error) {
	for _, rule := range f.rules {
		match, err := rule.Match(obj)
		if err != nil {
			return false, err
		}

		if !match {
			return false, nil
		}
	}

	return true, nil
}

func (f Filter) String() string {
	return fmt.Sprintf("Filter: %s", f.rules)
}
