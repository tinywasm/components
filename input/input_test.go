package input

import (
	"strings"
	"testing"
)

func TestInput_Render(t *testing.T) {
	i := &Input{
		Label:       "Username",
		Placeholder: "Enter username",
		Value:       "test",
		Type:        "text",
	}

	node := i.Render()

	// Input group container
	if node.Tag != "div" {
		t.Error("expected div tag for input group")
	}

	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "input-group") {
			hasClass = true
		}
	}
	if !hasClass {
		t.Error("expected input-group class")
	}
}
