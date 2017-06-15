package goui

import (
	"testing"
	"github.com/mooredwightd/gotestutil"
)

func TestGetUIConfig(t *testing.T) {
	c := GetUIConfig()
	gotestutil.AssertNotNil(t, c, "Expected non-nil result from GetUIConfig()")
	gotestutil.AssertGreaterThan(t, len(c.TemplatePaths()), 0, "Expected number of paths > 0.")
	t.Logf("TestGetUIConfig() TemplatePaths: %v\n", c.TemplatePaths())
}

func TestNewUIContext(t *testing.T) {
	c := NewUIContext()
	gotestutil.AssertNotNil(t, c, "Expected valid (non-nil) *UIContext.")
	gotestutil.AssertGreaterThan(t, len(c.p.GetStringSlice(CfgTemplatePath)), 0, "Expected number of paths > 0.")
	gotestutil.AssertNotNil(t, c.t, "Expected valid *Template field.")
	gotestutil.AssertTrue(t, c.p.GetBool(CfgVerbose), "Expected default value of Verbose of true.")
	gotestutil.AssertFalse(t, c.p.GetBool(CfgReload), "Expected default value of Reload of false.")
	gotestutil.AssertStringsEqual(t, c.p.GetString(CfgHomepage), "index.html",
		"Expected default value of home page. Actual: %s", "index.html")
	gotestutil.AssertStringsEqual(t, c.p.GetString(CfgPattern), "*.html",
		"Expected default value of search pattern. Actual: %s", "*.html")
}