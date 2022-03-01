package semver

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		text         string
		expected     Version
		errorMessage string
	}{
		{"0.0.1", Version{0, 0, 1}, ""},
		{"1.0.1", Version{1, 0, 1}, ""},
		{"-1.0.1", Version{0, 0, 0}, `Failed to parse version '-1.0.1': Invalid character(s) found in number "-1"`},
		{"01.0.1", Version{0, 0, 0}, `Failed to parse version '01.0.1': Version sections must not contain leading zeroes: "01"`},
		{"0.1", Version{0, 0, 0}, `Failed to parse version '0.1': Version text must be in 'Major.Minor.Patch' format`},
		{"", Version{0, 0, 0}, `Failed to parse version: Empty string`},
	}

	for _, tt := range tests {
		testname := tt.text
		t.Run(testname, func(t *testing.T) {
			v, err := Parse(tt.text)
			if v != tt.expected {
				t.Errorf("Expect to be '%v' but got '%v'", tt.expected, v)
			}
			if err != nil {
				if err.Error() != tt.errorMessage {
					t.Errorf("Expect error message '%s' but got '%s'.", tt.errorMessage, err.Error())
				}
			}
		})
	}
}

func TestLessThan(t *testing.T) {
	var tests = []struct {
		v1       Version
		v2       Version
		lessThan bool
	}{
		{
			Version{0, 1, 0},
			Version{1, 0, 0},
			true,
		},
		{
			Version{0, 1, 0},
			Version{0, 2, 0},
			true,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 2},
			true,
		},
		{
			Version{1, 1, 0},
			Version{0, 2, 0},
			false,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 1},
			false,
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s < %s: %v", tt.v1, tt.v2, tt.lessThan)
		t.Run(testName, func(t *testing.T) {
			res := tt.v1.LessThan(tt.v2)
			if res != tt.lessThan {
				t.Errorf("Expected to be '%v' but got '%v'", tt.lessThan, res)
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	var tests = []struct {
		v1          Version
		v2          Version
		greaterThan bool
	}{
		{
			Version{0, 1, 0},
			Version{1, 0, 0},
			false,
		},
		{
			Version{0, 1, 0},
			Version{0, 2, 0},
			false,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 2},
			false,
		},
		{
			Version{1, 1, 0},
			Version{0, 2, 0},
			true,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 1},
			false,
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s < %s: %v", tt.v1, tt.v2, tt.greaterThan)
		t.Run(testName, func(t *testing.T) {
			res := tt.v1.GreaterThan(tt.v2)
			if res != tt.greaterThan {
				t.Errorf("Expected to be '%v' but got '%v'", tt.greaterThan, res)
			}
		})
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		version Version
		text    string
	}{
		{Version{0, 0, 1}, "0.0.1"},
		{Version{1, 0, 0}, "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			if tt.version.String() != tt.text {
				t.Errorf("Expect to be '%s' but got '%s'.", tt.text, tt.version)
			}
		})
	}
}
