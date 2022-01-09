package app

import "testing"

func TestParse(t *testing.T) {
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

	app, err := Parse([]byte(input))

	if err != nil {
		t.Fatalf(`Expected to be succeed but failed`)
	}
	if app.Title != "Valid App 1" {
		t.Fatalf("failed")
	}
}
