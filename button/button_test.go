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

	node := btn.Render()

	if node.Tag != "button" {
		t.Error("expected button tag")
	}

	// Verify classes
	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "btn-primary") {
			hasClass = true
		}
	}

	if !hasClass {
		t.Error("expected btn-primary class")
	}

	// Verify text
	hasText := false
	for _, child := range node.Children {
		if s, ok := child.(string); ok && s == "Click" {
			hasText = true
		}
	}
	if !hasText {
		t.Error("expected text 'Click'")
	}
}
