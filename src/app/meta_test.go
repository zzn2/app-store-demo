package app

import (
	"testing"

	"github.com/zzn2/demo/appstore/semver"
)

var app = Meta{
	Title:   "Valid App 1",
	Version: semver.Version{Major: 0, Minor: 0, Patch: 1},
	Maintainers: []Maintainer{
		{
			Name:  "Alice",
			Email: "alice@hotmail.com",
		},
		{
			Name:  "Bob",
			Email: "bob@gmail.com",
		},
	},
	Company: "Random Inc.",
	Website: "https://website.com",
	Source:  "https://github.com/random/repo",
	License: "Apache-2.0",
	Description: `
 ### Intresting Title
 Some application content, and description
`,
}

func TestString(t *testing.T) {
	if app.String() != "App: Valid App 1@0.0.1" {
		t.Errorf("Failed")
	}
}
