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
func CreateRuleSet(queryParams map[string][]string, applyToObj interface{}) (RuleSet, error) {
	rs := RuleSet{}
	for key, value := range queryParams {
		if len(value) > 1 {
			return rs, errors.New(fmt.Sprintf("Key '%s' appeared multiple times with values of '%s'. Currently this case is not unsupported.", key, strings.Join(value, ", ")))
		}

		rule, err := NewRule(key, value[0], applyToObj)
		if err != nil {
			return rs, err
		}

		rs.AddRule(rule)
	}
	return rs, nil
}

// AddRule adds a new rule to the given RuleSet.
func (rs *RuleSet) AddRule(rule Rule) {
	rs.Rules = append(rs.Rules, rule)
}

// Match evaluates whether this Meta matches the given ruleset.
// It returns true if matches otherwise returns false.
// If any unexpected errors occurred during match operation, return the error.
func (rs RuleSet) Match(obj interface{}) (bool, error) {
	for _, rule := range rs.Rules {
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

func (rs RuleSet) String() string {
	return fmt.Sprintf("RuleSet: %s", rs.Rules)
}
