package filter

import (
	"testing"

	"github.com/zzn2/demo/appstore/app"
)

func TestMatch(t *testing.T) {
	flt, err := Create(map[string][]string{
		"title[like]": {"App"},
		"version":     {"0.0.1"},
	})
	if err != nil {
		t.Errorf("Failed to create filter")
	}

	input := `
title: Valid App 1
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`

	meta, err := app.Parse([]byte(input))
	if err != nil {
		t.Errorf("Failed to create app.Meta object.")
	}

	match, err := flt.Match(*meta)
	if err != nil {
		t.Errorf("Error occurred during Match operation: %s", err)
	}
	if !match {
		t.Errorf("Expected to match but failed.")
	}

}
