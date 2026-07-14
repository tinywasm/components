//go:build !wasm

package themetoggle

import (
	"testing"
)

func TestLabel_AllThemes_NonEmpty(t *testing.T) {
	themes := []TsTheme{TsThemeDark, TsThemeLight, "invalid"}
	for _, theme := range themes {
		l := label(theme)
		if l == "" {
			t.Errorf("label(%q) is empty", theme)
		}
	}
}

func TestToggle(t *testing.T) {
	tests := []struct {
		current TsTheme
		next    TsTheme
	}{
		{TsThemeDark, TsThemeLight},
		{TsThemeLight, TsThemeDark},
		{"invalid", TsThemeDark},
	}

	for _, tt := range tests {
		got := toggle(tt.current)
		if got != tt.next {
			t.Errorf("toggle(%q) = %q; want %q", tt.current, got, tt.next)
		}
	}
}

func TestValid(t *testing.T) {
	if !valid(TsThemeDark) {
		t.Error("valid(TsThemeDark) should be true")
	}
	if !valid(TsThemeLight) {
		t.Error("valid(TsThemeLight) should be true")
	}
	if valid("invalid") {
		t.Error("valid(\"invalid\") should be false")
	}
}

func TestIcon_AllThemes_NonEmpty(t *testing.T) {
	themes := []TsTheme{TsThemeDark, TsThemeLight, "invalid"}
	for _, theme := range themes {
		i := icon(theme)
		if i == "" {
			t.Errorf("icon(%q) is empty", theme)
		}
	}
}

func TestRenderCSS_NotEmpty(t *testing.T) {
	ts := &ThemeToggle{}
	css := ts.RenderCSS()
	if css.String() == "" {
		t.Error("RenderCSS() returned empty string")
	}
}
