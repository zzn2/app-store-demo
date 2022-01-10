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

func (s *Store) GetByTitle(title string) *Meta {
	return s.lastOrNil(func(app Meta) bool {
		return app.Title == title
	})
}

func (s *Store) GetByTitleAndVersion(title string, version string) *Meta {
	// If multiple found, return the last one, which is likely to be the latest version.
	// TODO: Should add version comparation logic here and only return the lastest version.
	return s.lastOrNil(func(app Meta) bool {
		return app.Title == title && app.Version == version
	})
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

func (s *Store) filter(match func(Meta) bool) []Meta {
	result := make([]Meta, 0)
	for _, app := range s.apps {
		if match(app) {
			result = append(result, app)
		}
	}

	return result
}

func (s *Store) firstOrNil(match func(Meta) bool) *Meta {
	for _, app := range s.apps {
		if match(app) {
			return &app
		}
	}

	return nil
}

func (s *Store) lastOrNil(match func(Meta) bool) *Meta {
	for i := len(s.apps) - 1; i >= 0; i-- {
		app := s.apps[i]
		if match(app) {
			return &app
		}
	}

	return nil
}
