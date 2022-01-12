// Package app provides functionalities of managing metadata of apps.
package app

import (
	"errors"
	"fmt"
	"sync"

	"github.com/zzn2/demo/appstore/filter"
)

// Store stores metadata of apps.
// They can be searched by various filters.
type Store struct {
	apps []Meta
}

var modifyLock sync.Mutex

// Add a new app metadata into the store.
// It returns error if the store already contains an app with the same title and version.
func (s *Store) Add(app Meta) error {
	if s.GetByTitleAndVersion(app.Title, app.Version) != nil {
		return errors.New(fmt.Sprintf("App '%s' with version '%s' already exists.", app.Title, app.Version))
	}

	modifyLock.Lock()
	defer modifyLock.Unlock()
	// add lock to append operation to avoid potential racing cases.
	s.apps = append(s.apps, app)
	return nil
}

// GetByTitle gets an app metadata using title.
// It returns the matching metadata if exists, otherwise returns nil.
// If multiple version exists for the same title, it returns the last saved version.
func (s *Store) GetByTitle(title string) *Meta {
	return s.lastOrNil(func(app Meta) bool {
		return app.Title == title
	})
}

// GetByTitleAndVersion gets an app metadata using title and version.
// It returns the matching metadata if exists, otherwise returns nil.
func (s *Store) GetByTitleAndVersion(title string, version string) *Meta {
	// If multiple found, return the last one, which is likely to be the latest version.
	// TODO: Should add version comparing logic here and only return the latest version.
	return s.lastOrNil(func(app Meta) bool {
		return app.Title == title && app.Version == version
	})
}

// List returns the list of stored apps.
// It accepts a filter.RuleSet as parameter, only apps matching the rules could be listed in the result.
// RuleSet could also be empty (i.e. contains no rules). In this case, all the apps will be listed in the result.
// If no matching apps found, return an empty slice.
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

// filter returns the apps that matches the given rule in the store.
// It returns an empty slice if no matching apps found.
func (s *Store) filter(match func(Meta) bool) []Meta {
	result := make([]Meta, 0)
	for _, app := range s.apps {
		if match(app) {
			result = append(result, app)
		}
	}

	return result
}

// firstOrNil gets the first app that matches the given rule.
// It returns nil if no matching apps were found.
func (s *Store) firstOrNil(match func(Meta) bool) *Meta {
	for _, app := range s.apps {
		if match(app) {
			return &app
		}
	}

	return nil
}

// lastOrNil gets the last app that matches the given rule.
// It returns nil if no matching apps were found.
func (s *Store) lastOrNil(match func(Meta) bool) *Meta {
	for i := len(s.apps) - 1; i >= 0; i-- {
		app := s.apps[i]
		if match(app) {
			return &app
		}
	}

	return nil
}
