package app

import "gopkg.in/yaml.v2"

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

func (m *Meta) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}
