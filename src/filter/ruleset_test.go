package filter

import (
	"testing"
)

func TestCreate(t *testing.T) {
	ruleSet, err := CreateRuleSet(map[string][]string{
		"title[like]": {"App"},
		"version":     {"0.0.1"},
	})
	if err != nil {
		t.Errorf("Failed to create RuleSet: %e", err)
	}

	if len(ruleSet.Rules) != 2 {
		t.Errorf("Expected to be 2 tules but got %d", len(ruleSet.Rules))
	}
}

func TestCreate_DuplicateKey(t *testing.T) {
	ruleSet, err := CreateRuleSet(map[string][]string{
		"title": {"App1", "App2"},
	})
	if err == nil {
		t.Errorf("Expected to have error but had none.")
	}

	expectedErrMsg := "Key 'title' appeared multiple times with values of 'App1, App2'. Currently this case is not unsupported."
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected to have error '%s' but got '%s'", expectedErrMsg, err.Error())
	}

	if ruleSet != nil {
		t.Errorf("Expected ruleSet to be nil but not nil.")
	}
}

func TestCreate_RuleCreationFailure(t *testing.T) {
	ruleSet, err := CreateRuleSet(map[string][]string{
		"title[dummy]": {"App1"},
	})
	if err == nil {
		t.Errorf("Expected to have error but had none.")
	}

	expectedErrMsg := "Unrecognized operator type 'dummy'"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected to have error '%s' but got '%s'", expectedErrMsg, err.Error())
	}

	if ruleSet != nil {
		t.Errorf("Expected ruleSet to be nil but not nil.")
	}
}
