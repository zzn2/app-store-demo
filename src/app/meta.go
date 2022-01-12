package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zzn2/demo/appstore/filter"
	"gopkg.in/yaml.v2"
)

type Maintainer struct {
	Name  string `binding:"required"`
	Email string `binding:"required,email"`
}

type Meta struct {
	Title       string       `binding:"required"`
	Version     string       `binding:"required"`
	Maintainers []Maintainer `binding:"required,dive"`
	Company     string       `binding:"required"`
	Website     string       `binding:"required,url"`
	Source      string       `binding:"required"`
	License     string       `binding:"required"`
	Description string       `binding:"required"`
}

// Parse parses an yaml text into the Meta struct.
// It returns nil along with the error if any error occurred during the parse.
func Parse(data []byte) (*Meta, error) {
	var m Meta
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// MatchRule evaluates whether this Meta matches the given rule.
// It returns true if matches otherwise returns false.
// If any unexpected errors occurred during match operation, return the error.
func (m *Meta) MatchRule(rule filter.Rule) (bool, error) {
	// TODO: embed the logic into the struct using struct tags
	// https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go
	switch strings.ToLower(rule.FieldName) {
	case "title":
		return rule.Evaluate(m.Title)
	case "version":
		return rule.Evaluate(m.Version)
	case "maintainer.name":
		for _, maintainer := range m.Maintainers {
			match, err := rule.Evaluate(maintainer.Name)
			if err != nil {
				return false, err
			}
			if match {
				// Early return when the first match found.
				// If not found, do not return and go to next loop.
				return true, nil
			}
		}
		// All maintainers not match if goes to this line.
		return false, nil
	case "company":
		return rule.Evaluate(m.Company)
	case "description":
		return rule.Evaluate(m.Description)
	default:
		return false, errors.New(fmt.Sprintf("Unsupported rule for field '%s'", rule.FieldName))
	}
}

// MatchRule evaluates whether this Meta matches the given ruleset.
// It returns true if matches otherwise returns false.
// If any unexpected errors occurred during match operation, return the error.
func (m *Meta) MatchRuleSet(ruleSet filter.RuleSet) (bool, error) {
	for _, rule := range ruleSet.Rules {
		match, err := m.MatchRule(rule)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

// String returns the string representation of this object.
func (m Meta) String() string {
	return fmt.Sprintf("App: %s@%s", m.Title, m.Version)
}
