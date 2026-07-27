//go:build !wasm

package fieldset

import "testing"

// TestRenderCSS_StylesFieldset guards that Style outputs proper rules for root and parts.
func TestRenderCSS_StylesFieldset(t *testing.T) {
	cssStr := (&Fieldset{}).Style().Stylesheet().String()
	if cssStr == "" {
		t.Fatal("Stylesheet() returned empty")
	}
	for _, want := range []string{".fieldset", ".fieldset__label"} {
		if !contains(cssStr, want) {
			t.Errorf("Stylesheet() missing rule for %q", want)
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
