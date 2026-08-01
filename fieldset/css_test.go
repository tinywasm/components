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
	// (no vertical padding — ChipBox alone) and packs its text at the leading
	// edge. Inline padding is the one exemption: it leaves the height alone
	// and keeps the text off the chip's filled edges.
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
	if contains(b, "padding: ") || contains(b, "padding-block") {
		t.Errorf("label chip must not grow vertically (keeps it at badge height), block:\n%s", b)
	}
	if !contains(b, "padding-inline") {
		t.Errorf("label chip text needs inline air off the chip's edges, block:\n%s", b)
	}
	if !contains(b, "justify-content: flex-start") {
		t.Errorf("label chip must pack text at the leading edge, block:\n%s", b)
	}
}

func TestFieldset_ErrorClearsTheInputEdges(t *testing.T) {
	// The error message sits inside the input's top-right with matching air on
	// both edges. An absolute box lays out against the Root's PADDING box, so a
	// Space2 pin only buys back the Root's own Pad and left the text flush with
	// the input's border; and straddling the border cut the text in half.
	css := (&Fieldset{}).RenderCSS().String()
	i := strings.Index(css, ".tw-field__error {")
	if i == -1 {
		t.Fatal("expected a rule for .tw-field__error")
	}
	body := css[i:]
	end := strings.Index(body, "}")
	if end == -1 {
		t.Fatal("malformed rule block")
	}
	b := body[:end]
	for _, want := range []string{
		"inset-inline-end: var(--space-4,",
		"inset-block-start: var(--space-4,",
	} {
		if !contains(b, want) {
			t.Errorf("error text needs %q to clear the input's border, block:\n%s", want, b)
		}
	}
	if contains(b, "transform: translateY(") {
		t.Errorf("error text must sit inside the box, not straddle the input's top line, block:\n%s", b)
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
