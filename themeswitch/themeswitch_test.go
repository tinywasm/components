package themeswitch

import (
	"testing"
)

func TestLabel_AllThemes_NonEmpty(t *testing.T) {
	themes := []Theme{ThemeAuto, ThemeDark, ThemeLight, "invalid"}
	for _, theme := range themes {
		l := label(theme)
		if l == "" {
			t.Errorf("label(%q) is empty", theme)
		}
	}
}

func TestCycle(t *testing.T) {
	tests := []struct {
		current Theme
		next    Theme
	}{
		{ThemeAuto, ThemeDark},
		{ThemeDark, ThemeLight},
		{ThemeLight, ThemeAuto},
		{"invalid", ThemeDark},
	}

	for _, tt := range tests {
		got := cycle(tt.current)
		if got != tt.next {
			t.Errorf("cycle(%q) = %q; want %q", tt.current, got, tt.next)
		}
	}
}

func TestValid(t *testing.T) {
	if !valid(ThemeAuto) {
		t.Error("valid(ThemeAuto) should be true")
	}
	if !valid(ThemeDark) {
		t.Error("valid(ThemeDark) should be true")
	}
	if !valid(ThemeLight) {
		t.Error("valid(ThemeLight) should be true")
	}
	if valid("invalid") {
		t.Error("valid(\"invalid\") should be false")
	}
}

func TestRenderCSS_NotEmpty(t *testing.T) {
	ts := &ThemeSwitch{}
	css := ts.RenderCSS()
	if css == "" {
		t.Error("RenderCSS() returned empty string")
	}
}
