package card

import (
	"strings"
	"testing"

	dom "github.com/tinywasm/components/internal/dom"
)

func TestCard_Render(t *testing.T) {
	c := &Card{
		Header: dom.Tag("h3").Text("Header").ToComponent(),
		Body:   dom.Tag("p").Text("Body").ToComponent(),
		Footer: dom.Tag("div").Text("Footer").ToComponent(),
	}

	node := c.Render()

	if node.Tag != "div" {
		t.Error("expected div tag")
	}

	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "card") {
			hasClass = true
		}
	}
	if !hasClass {
		t.Error("expected card class")
	}

	// Check for children (header, body, footer)
	// This is hard to check deeply because children are any, and might be nested.
	// But we can check count or structure roughly.
	if len(node.Children) < 3 {
		t.Error("expected at least 3 children")
	}
}
