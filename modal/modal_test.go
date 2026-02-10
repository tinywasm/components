package modal

import (
	"strings"
	"testing"
)

func TestModal_Render(t *testing.T) {
	m := &Modal{
		Title:   "Example Modal",
		Content: nil, // Can be nil
		Visible: false,
	}

	node := m.Render()

	if node.Tag != "div" {
		t.Error("expected div tag")
	}

	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "modal") {
			hasClass = true
		}
	}
	if !hasClass {
		t.Error("expected modal class")
	}

	// Should be hidden by default
	hasHidden := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "hidden") {
			hasHidden = true
		}
	}
	if !hasHidden {
		t.Error("expected hidden class")
	}
}
