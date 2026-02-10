package nav

import (
	"strings"
	"testing"
)

func TestNav_Render(t *testing.T) {
	n := &Nav{
		Items: []NavItem{
			{Label: "Home", Route: "home"},
			{Label: "Profile", Route: "profile"},
		},
	}

	node := n.Render()

	if node.Tag != "nav" {
		t.Error("expected nav tag")
	}

	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "nav") {
			hasClass = true
		}
	}
	if !hasClass {
		t.Error("expected nav class")
	}
}
