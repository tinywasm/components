package table

import (
	"strings"
	"testing"
)

func TestTable_Render(t *testing.T) {
	tbl := &Table{
		Headers: []string{"Name", "Age"},
		Rows: [][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	}

	html := tbl.Render().RenderHTML()

	if !strings.HasPrefix(html, "<table") {
		t.Error("expected table tag")
	}
	if !strings.Contains(html, "class='table'") {
		t.Error("expected table class")
	}

	// Check headers
	if !strings.Contains(html, "<th>Name</th>") {
		t.Error("expected Name header")
	}
	if !strings.Contains(html, "<th>Age</th>") {
		t.Error("expected Age header")
	}

	// Check rows
	if !strings.Contains(html, "<td>Alice</td>") {
		t.Error("expected Alice cell")
	}
	if !strings.Contains(html, "<td>30</td>") {
		t.Error("expected 30 cell")
	}
}
