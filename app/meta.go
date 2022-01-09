package app

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type Meta struct {
	Title       string
	Version     string
	Maintainers []struct {
		Name  string
		Email string
	}
	Company     string
	Website     string
	Source      string
	License     string
	Description string
}

func Parse(data []byte) (*Meta, error) {
	var m Meta
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m Meta) String() string {
	return fmt.Sprintf("App: %s@%s", m.Title, m.Version)
}
