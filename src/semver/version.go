// Package semver provides parsing and validation logic for [Semantic Version](https://semver.org/)s.
//
// Note:
//   1. Here we only implement subset of Semantic Versioning just for demo usage.
//   2. Most of the code referred the implemention from https://github.com/blang/semver.
package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

var Empty = Version{}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}

func parseSection(text string) (uint64, error) {
	if !containsOnly(text, "0123456789") {
		return 0, fmt.Errorf("Invalid character(s) found in number %q", text)
	}
	if hasLeadingZeroes(text) {
		return 0, fmt.Errorf("Version sections must not contain leading zeroes: %q", text)
	}

	res, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// Parse parses a text into a Version object.
// errors will be returned if the text is not a valid format.
func Parse(s string) (Version, error) {
	if len(s) == 0 {
		return Version{}, errors.New("Failed to parse version: Empty string")
	}

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("Failed to parse version '%s': Version text must be in 'Major.Minor.Patch' format", s)
	}

	sections := make([]uint64, 3)
	for i := range sections {
		val, err := parseSection(parts[i])
		if err != nil {
			return Version{}, fmt.Errorf("Failed to parse version '%s': %w", s, err)
		}
		sections[i] = val
	}

	v := Version{}
	v.Major = sections[0]
	v.Minor = sections[1]
	v.Patch = sections[2]

	return v, nil
}

// UnmarshalYAML will be called when deserializing a Version object from part of YAML text.
func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var text string
	err := unmarshal(&text)
	if err != nil {
		return err
	}

	version, err := Parse(text)
	if err != nil {
		return err
	}

	v.Major = version.Major
	v.Minor = version.Minor
	v.Patch = version.Patch
	return nil
}

// MarshalJSON will be called when serializing a Version object into text as part of json.
func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", v.String())), nil
}

// String returns the string representation of the Version object.
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
