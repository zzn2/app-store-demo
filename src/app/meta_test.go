package app

import (
	"testing"

	"example.com/goserver/filter"
)

const sampleInput = `
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

func TestParse(t *testing.T) {
	app, err := Parse([]byte(sampleInput))

	if err != nil {
		t.Fatalf(`Expected to be succeed but failed`)
	}
	if app.Title != "Valid App 1" {
		t.Fatalf("failed")
	}
}

func TestMatchRule(t *testing.T) {
	var tests = []struct {
		ruleText       string
		expectedResult bool
		expectedErrMsg string
	}{
		{
			"title[like]=App",
			true,
			"",
		},
		{
			"title=App",
			false,
			"",
		},
		{
			"title=Valid App 1",
			true,
			"",
		},
		{
			"maintainer.name[like]=second",
			true,
			"",
		},
		{
			"maintainer.name[like]=third",
			false,
			"",
		},
	}

	app, err := Parse([]byte(sampleInput))
	if err != nil {
		t.Fatalf("Failed to parse app meta.")
	}

	for _, tt := range tests {
		testName := tt.ruleText
		t.Run(testName, func(t *testing.T) {
			rule, err := filter.ParseRule(tt.ruleText)
			if err != nil {
				t.Fatalf("Failed to parse rule %s", tt.ruleText)
			}
			match, err := app.MatchRule(*rule)
			if match != tt.expectedResult {
				t.Errorf("Expect '%v' but got '%v'.", tt.expectedResult, match)
			}
			if err != nil {
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expect err to be '%s' but got '%s'.", tt.expectedErrMsg, err.Error())
				}
			}
		})
	}
}
