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
	tbl.Init(nil)

	html := tbl.Render().String()

	if !HasPrefix(html, "<table") {
		t.Error("expected table tag")
	}
	if !Contains(html, "class='datatable'") {
		t.Error("expected datatable class")
	}

	// Check headers
	if !Contains(html, "role='columnheader'>Name</th>") {
		t.Error("expected Name header")
	}
	if !Contains(html, "role='columnheader'>Age</th>") {
		t.Error("expected Age header")
	}

	// Check rows
	if !Contains(html, "role='gridcell'>Alice</td>") {
		t.Error("expected Alice cell")
	}
	if !Contains(html, "role='gridcell'>30</td>") {
		t.Error("expected 30 cell")
	}
}
