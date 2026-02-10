package table

import (
	"strings"
	"testing"
)

func TestTable_Render(t *testing.T) {
	tab := &Table{
		Headers: []string{"ID", "Name"},
		Rows: [][]string{
			{"1", "Alice"},
			{"2", "Bob"},
		},
	}

	node := tab.Render()

	if node.Tag != "table" {
		t.Error("expected table tag")
	}

	hasClass := false
	for _, attr := range node.Attrs {
		if attr.Key == "class" && strings.Contains(attr.Value, "table") {
			hasClass = true
		}
	}
	if !hasClass {
		t.Error("expected table class")
	}
}
