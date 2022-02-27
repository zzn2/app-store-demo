// Package semver provides parsing and validation logic for [Semantic Versioning](https://semver.org/).
//
// Note:
//   1. Here we only implementing a subset of Semantic Versioning just for demo usage.
//   2. Most of the code referred to the implemention of https://github.com/blang/semver.
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

func Parse(s string) (Version, error) {
	if len(s) == 0 {
		return Version{}, errors.New("Version string empty")
	}

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return Version{}, errors.New("No Major.Minor.Patch elements found")
	}

	major, err := parseSection(parts[0])
	if err != nil {
		return Version{}, err
	}

	minor, err := parseSection(parts[1])
	if err != nil {
		return Version{}, err
	}

	patch, err := parseSection(parts[2])
	if err != nil {
		return Version{}, err
	}

	v := Version{}
	v.Major = major
	v.Minor = minor
	v.Patch = patch

	return v, nil
}

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

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", v.String())), nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
