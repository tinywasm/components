//go:build !wasm

package fieldset

import (
	"regexp"
	"strings"
	"testing"
)

var (
	classRegex = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
)

// TestRenderCSS_StylesFieldset guards that Style outputs proper rules for root and parts.
func TestRenderCSS_StylesFieldset(t *testing.T) {
	cssStr := (&Fieldset{}).RenderCSS().String()
	if cssStr == "" {
		t.Fatal("Stylesheet() returned empty")
	}
	for _, want := range []string{".tw-field", ".tw-field__label"} {
		if !contains(cssStr, want) {
			t.Errorf("Stylesheet() missing rule for %q", want)
		}
	}
}

func TestFieldset_LabelChip(t *testing.T) {
	// PLAN v0.2.0 item 4: the label chip matches the list badge's height
	// (no Pad — ChipBox alone) and packs its text at the leading edge.
	css := (&Fieldset{}).RenderCSS().String()
	i := strings.Index(css, ".tw-field__label {")
	if i == -1 {
		t.Fatal("expected a rule for .tw-field__label")
	}
	body := css[i:]
	end := strings.Index(body, "}")
	if end == -1 {
		t.Fatal("malformed rule block")
	}
	b := body[:end]
	if contains(b, "padding") {
		t.Errorf("label chip must not be padded (keeps it at badge height), block:\n%s", b)
	}
	if !contains(b, "justify-content: flex-start") {
		t.Errorf("label chip must pack text at the leading edge, block:\n%s", b)
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

	filterClasses := func(classes map[string]bool, prefix string) map[string]bool {
		filtered := make(map[string]bool)
		for cls := range classes {
			if strings.HasPrefix(cls, prefix) {
				filtered[cls] = true
			}
		}
		return filtered
	}

	f := &Fieldset{}
	css := f.RenderCSS().String()
	cssClasses := filterClasses(extractCSSClasses(css), "tw-field")

	expectedFormClasses := map[string]bool{
		"tw-field":              true,
		"tw-field__label":       true,
		"tw-field__input":       true,
		"tw-field__error":       true,
		"tw-field__radio-group": true,
	}

	for cls := range cssClasses {
		if !expectedFormClasses[cls] {
			t.Errorf("Fieldset CSS contains unexpected class %q which form v0.3.0 does not emit", cls)
		}
	}
	for cls := range expectedFormClasses {
		if !cssClasses[cls] {
			t.Errorf("Fieldset CSS missing style rule for form v0.3.0 class %q", cls)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
