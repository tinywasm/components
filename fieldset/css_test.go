//go:build !wasm

package fieldset

import "testing"

// The skin exists solely to emit CSS for the form structure. Guard that it stays
// non-empty and keeps targeting the .tw-field wrapper + its label.
func TestRenderCSS_StylesTwField(t *testing.T) {
	css := (&Fieldset{}).RenderCSS().String()
	if css == "" {
		t.Fatal("RenderCSS() returned empty")
	}
	for _, want := range []string{".tw-field", ".tw-field label"} {
		if !contains(css, want) {
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
