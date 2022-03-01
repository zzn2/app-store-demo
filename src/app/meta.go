package app

import (
	"fmt"

	"github.com/zzn2/demo/appstore/semver"
)

type Maintainer struct {
	Name  string `binding:"required"`
	Email string `binding:"required,email"`
}

type Meta struct {
	Title       string         `binding:"required"`
	Version     semver.Version `binding:"required"`
	Maintainers []Maintainer   `binding:"required,dive"`
	Company     string         `binding:"required"`
	Website     string         `binding:"required,url"`
	Source      string         `binding:"required"`
	License     string         `binding:"required"`
	Description string         `binding:"required"`
}

// String returns the string representation of this object.
func (m Meta) String() string {
	return fmt.Sprintf("App: %s@%s", m.Title, m.Version)
}
