package form

import (
	"strings"
	"testing"
)

func TestForm_Render(t *testing.T) {
	f := &Form{
		Fields: nil,
	}

	node := f.Render()

	if node.Tag != "form" {
		t.Error("expected form tag")
	}

	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "form") {
			hasClass = true
		}
	}
	if !hasClass {
		t.Error("expected form class")
	}
}
