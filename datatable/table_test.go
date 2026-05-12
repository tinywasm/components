package datatable

import (
	"testing"

	. "github.com/tinywasm/fmt"
)

func TestTable_Render(t *testing.T) {
	tbl := &DataTable{
		Headers: []string{"Name", "Age"},
		Rows: [][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	}

	html := tbl.Render().RenderHTML()

	if !HasPrefix(html, "<table") {
		t.Error("expected table tag")
	}
	if !Contains(html, "class='table'") {
		t.Error("expected table class")
	}

	// Check headers
	if !Contains(html, "<th>Name</th>") {
		t.Error("expected Name header")
	}
	if !Contains(html, "<th>Age</th>") {
		t.Error("expected Age header")
	}

	// Check rows
	if !Contains(html, "<td>Alice</td>") {
		t.Error("expected Alice cell")
	}
	if !Contains(html, "<td>30</td>") {
		t.Error("expected 30 cell")
	}
}
