//go:build !wasm

package themetoggle

import (
	"regexp"
	"strings"
	"testing"
)

var (
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
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
	css := ts.Style().Stylesheet()
	if css.String() == "" {
		t.Error("Stylesheet() returned empty string")
	}
}

func TestPairMarkupAndStylesheet(t *testing.T) {
	extractCSSClasses := func(css string) map[string]bool {
		classes := make(map[string]bool)
		matches := classRegex.FindAllStringSubmatch(css, -1)
		for _, m := range matches {
			if len(m) > 1 {
				classes[m[1]] = true
			}
		}
		return classes
	}

	extractHTMLClasses := func(html string) map[string]bool {
		classes := make(map[string]bool)
		matches1 := htmlClassRegex.FindAllStringSubmatch(html, -1)
		for _, m := range matches1 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		matches2 := htmlClassRegex2.FindAllStringSubmatch(html, -1)
		for _, m := range matches2 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		return classes
	}

	filterClasses := func(classes map[string]bool, prefix string) map[string]bool {
		filtered := make(map[string]bool)
		for cls := range classes {
			if strings.HasPrefix(cls, prefix) {
				filtered[cls] = true
			}
		}
		return filtered
	}

	tt := &ThemeToggle{}
	tt.Init(nil)
	html := tt.Render().String()
	css := tt.Style().Stylesheet().String()

	htmlClasses := filterClasses(extractHTMLClasses(html), "themetoggle")
	cssClasses := filterClasses(extractCSSClasses(css), "themetoggle")

	for cls := range cssClasses {
		if !htmlClasses[cls] {
			t.Errorf("ThemeToggle CSS class %q does not exist in rendered HTML", cls)
		}
	}
	for cls := range htmlClasses {
		if !cssClasses[cls] {
			t.Errorf("ThemeToggle HTML class %q is unstyled in CSS", cls)
		}
	}
}
