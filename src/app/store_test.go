package app

import (
	"testing"

	"github.com/zzn2/demo/appstore/filter"
	"github.com/zzn2/demo/appstore/semver"
)

var (
	v_0_0_1 = semver.Version{Major: 0, Minor: 0, Patch: 1}
	v_0_0_2 = semver.Version{Major: 0, Minor: 0, Patch: 2}
	v_0_0_3 = semver.Version{Major: 0, Minor: 0, Patch: 3}
)

var (
	app1v1 = Meta{
		Title:   "App1",
		Version: v_0_0_1,
	}

	app1v2 = Meta{
		Title:   "App1",
		Version: v_0_0_2,
	}

	app2v1 = Meta{
		Title:   "App2",
		Version: v_0_0_1,
	}
)

func TestAdd(t *testing.T) {
	var store Store

	if len(store.apps) != 0 {
		t.Errorf("Expected store is empty but contained %d apps.", len(store.apps))
	}

	err := store.Add(app1v1)

	if len(store.apps) != 1 {
		t.Errorf("Expected store contains 1 app but actually contained %d apps.", len(store.apps))
	}
	if err != nil {
		t.Errorf("Should not have error but error '%s' occurred.", err.Error())
	}

	err = store.Add(app1v2)

	if len(store.apps) != 2 {
		t.Errorf("Expected store contains 2 apps but actually contained %d apps.", len(store.apps))
	}
	if err != nil {
		t.Errorf("Should not have error but error '%s' occurred.", err.Error())
	}

	err = store.Add(app1v2)

	if len(store.apps) != 2 {
		t.Errorf("Expected store contains 2 apps but actually contained %d apps.", len(store.apps))
	}
	if err == nil {
		expectedErrMsg := "App 'App1' with version '0.0.2' already exists."
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message to be '%s' but got '%s'.", expectedErrMsg, err.Error())
		}
	}
}

func equals(app1 Meta, app2 Meta) bool {
	return app1.Title == app2.Title && app1.Version == app2.Version
}

func TestGetByTitle(t *testing.T) {
	var store Store

	store.Add(app1v1)
	store.Add(app1v2)
	store.Add(app2v1)

	result1 := store.GetByTitle("App1")
	if !equals(*result1, app1v2) {
		t.Errorf("Expected to be '%s' but got '%s'", result1, app1v2)
	}

	result2 := store.GetByTitle("App2")
	if !equals(*result2, app2v1) {
		t.Errorf("Expected to be '%s' but got '%s'", result2, app2v1)
	}

	result3 := store.GetByTitle("App3")
	if result3 != nil {
		t.Errorf("Expected to be nil but got '%s'", result3)
	}
}

func TestGetByTitleAndVersion(t *testing.T) {
	var store Store

	store.Add(app1v1)
	store.Add(app1v2)
	store.Add(app2v1)

	result1 := store.GetByTitleAndVersion("App1", v_0_0_1)
	if !equals(*result1, app1v1) {
		t.Errorf("Expected to be '%s' but got '%s'", result1, app1v2)
	}

	result2 := store.GetByTitleAndVersion("App1", v_0_0_2)
	if !equals(*result2, app1v2) {
		t.Errorf("Expected to be '%s' but got '%s'", result2, app1v2)
	}

	result3 := store.GetByTitleAndVersion("App1", v_0_0_3)
	if result3 != nil {
		t.Errorf("Expected to be nil but got '%s'", result3)
	}

	result4 := store.GetByTitleAndVersion("App2", v_0_0_1)
	if !equals(*result4, app2v1) {
		t.Errorf("Expected to be '%s' but got '%s'", result4, app2v1)
	}

	result5 := store.GetByTitleAndVersion("App3", v_0_0_1)
	if result5 != nil {
		t.Errorf("Expected to be nil but got '%s'", result5)
	}
}

func TestList(t *testing.T) {
	var store Store

	store.Add(app1v1)
	store.Add(app1v2)
	store.Add(app2v1)

	var ruleSet filter.RuleSet

	result, err := store.List(ruleSet)
	if err != nil {
		t.Errorf("Expected to be no error but got '%s'", err.Error())
	}
	if len(result) != 3 {
		t.Errorf("Expected to be 3 items but got %d", len(result))
	}

	rule, _ := filter.ParseRule("title=App1")
	ruleSet.AddRule(*rule)

	result, err = store.List(ruleSet)
	if err != nil {
		t.Errorf("Expected to be no error but got '%s'", err.Error())
	}
	if len(result) != 2 {
		t.Errorf("Expected to be 2 items but got %d", len(result))
	}

	rule, _ = filter.ParseRule("unknown=App1")
	ruleSet.AddRule(*rule)

	result, err = store.List(ruleSet)
	if err != nil {
		expectedErrMsg := "Error occurred during searching app: Unsupported rule for field 'unknown'"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message to be '%s' but got '%s'", expectedErrMsg, err.Error())
		}
	}
	if len(result) != 0 {
		t.Errorf("Expected to be 0 items but got %d", len(result))
	}
}
