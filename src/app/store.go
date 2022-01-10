package app

import (
	"errors"
	"fmt"

	"github.com/zzn2/demo/appstore/filter"
)

type Store struct {
	apps []Meta
}

func (s *Store) Add(app Meta) {
	s.apps = append(s.apps, app)
}

func (s *Store) List(ruleSet filter.RuleSet) ([]Meta, error) {
	result := make([]Meta, 0)
	for _, app := range s.apps {
		matched, err := app.MatchRuleSet(ruleSet)
		if err != nil {
			return result, errors.New(fmt.Sprintf("Error occurred during searching app: %s", err))
		}

		if matched {
			result = append(result, app)
		}
	}

	return result, nil
}
