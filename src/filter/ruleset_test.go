package filter

import (
	"testing"

	"github.com/zzn2/demo/appstore/semver"
)

type Meta struct {
	Title   string
	Version semver.Version
}

var app Meta

func TestCreate(t *testing.T) {
	ruleSet, err := CreateRuleSet(map[string][]string{
		"title[like]": {"App"},
		"version":     {"0.0.1"},
	}, app)
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
	}, app)
	if err == nil {
		t.Errorf("Expected to have error but had none.")
	}

	expectedErrMsg := "Key 'title' appeared multiple times with values of 'App1, App2'. Currently this case is not unsupported."
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected to have error '%s' but got '%s'", expectedErrMsg, err.Error())
	}

	if len(ruleSet.Rules) != 0 {
		t.Errorf("Expected ruleSet to be empty but not empty.")
	}
}

func TestCreate_RuleCreationFailure(t *testing.T) {
	ruleSet, err := CreateRuleSet(map[string][]string{
		"title[dummy]": {"App1"},
	}, app)
	if err == nil {
		t.Errorf("Expected to have error but had none.")
	}

	expectedErrMsg := "Unrecognized operator type 'dummy'"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected to have error '%s' but got '%s'", expectedErrMsg, err.Error())
	}

	if len(ruleSet.Rules) != 0 {
		t.Errorf("Expected ruleSet to be empty but not empty.")
	}
}
