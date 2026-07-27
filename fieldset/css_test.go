//go:build !wasm

package fieldset

import "testing"

// TestRenderCSS_StylesFieldset guards that RenderCSS outputs proper rules for root and parts.
func TestRenderCSS_StylesFieldset(t *testing.T) {
	cssStr := (&Fieldset{}).RenderCSS().String()
	if cssStr == "" {
		t.Fatal("RenderCSS() returned empty")
	}
	for _, want := range []string{".fieldset", ".fieldset__label"} {
		if !contains(cssStr, want) {
			t.Errorf("RenderCSS() missing rule for %q", want)
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
