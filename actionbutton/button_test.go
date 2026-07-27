package actionbutton

import (
	"strings"
	"testing"
)

func TestButton_Render(t *testing.T) {
	btn := &ActionButton{
		Text:    "Click",
		Variant: "primary",
	}

	html := btn.Render().String()

	if !strings.HasPrefix(html, "<button") {
		t.Error("expected button tag, got: " + html)
	}

	// Verify classes
	if !strings.Contains(html, "actionbutton") || !strings.Contains(html, "actionbutton__primary") {
		t.Error("expected actionbutton classes, got: " + html)
	}

	// Verify text
	if !strings.Contains(html, ">Click<") {
		t.Error("expected text 'Click', got: " + html)
	}
}
