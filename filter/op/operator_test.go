package op

import "testing"

func TestParse(t *testing.T) {
	var tests = []struct {
		input        string
		want         Operator
		errorMessage string
	}{
		{"", Equals, ""},
		{"like", Like, ""},
		{"Like", Like, ""},
		{"LIKE", Like, ""},
		{"lt", LessThan, ""},
		{"gt", GreaterThan, ""},
		{"other", Unknown, "Unrecognized operation type 'other'"},
	}

	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			op, err := Parse(tt.input)
			if op != tt.want {
				t.Errorf("Expect '%s' but got '%s'.", tt.want, op)
			}
			if err != nil {
				if err.Error() != tt.errorMessage {
					t.Errorf("Expect error message '%s' but got '%s'.", tt.errorMessage, err.Error())
				}
			}
		})
	}
}
