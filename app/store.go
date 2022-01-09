package app

import (
	"errors"
	"fmt"

	"example.com/goserver/filter"
)

type Store struct {
	apps []Meta
}

func (s *Store) Add(app Meta) {
	s.apps = append(s.apps, app)
}

func (s *Store) List(flt filter.Filter) ([]Meta, error) {
	result := make([]Meta, 0)
	for _, app := range s.apps {
		matched, err := flt.Match(app)
		if err != nil {
			return result, errors.New(fmt.Sprintf("Error occurred during searching app: %s", err))
		}

		if matched {
			result = append(result, app)
		}
	}

	return result, nil
}
