package button

import (
	"strings"
	"testing"
)

func TestButton_Render(t *testing.T) {
	btn := &Button{
		Text:    "Click",
		Variant: "primary",
	}

	html := btn.Render().RenderHTML()

	if !strings.HasPrefix(html, "<button") {
		t.Error("expected button tag, got: " + html)
	}

	// Verify classes
	if !strings.Contains(html, "class='btn btn-primary'") && !strings.Contains(html, "class='btn-primary btn'") {
		t.Error("expected btn-primary class, got: " + html)
	}

	// Verify text
	if !strings.Contains(html, ">Click<") {
		t.Error("expected text 'Click', got: " + html)
	}
}
