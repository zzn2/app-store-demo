package semver

import "testing"

func TestEvaluate(t *testing.T) {
	var tests = []struct {
		text         string
		expected     Version
		errorMessage string
	}{
		{"0.0.1", Version{0, 0, 1}, ""},
		{"1.0.1", Version{1, 0, 1}, ""},
		{"-1.0.1", Version{0, 0, 0}, `Invalid character(s) found in number "-1"`},
		{"01.0.1", Version{0, 0, 0}, `Version sections must not contain leading zeroes: "01"`},
		{"0.1", Version{0, 0, 0}, `No Major.Minor.Patch elements found`},
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
