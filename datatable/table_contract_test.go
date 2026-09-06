package datatable

import (
	"testing"

	"webtyp.com/dom"
)

var _ dom.Component = (*DataTable)(nil)

func TestDataTable_InitTwiceSafe(t *testing.T) {
	tbl := &DataTable{Headers: []string{"A"}, Rows: [][]string{{"1"}}}

	tbl.Init(nil)
	tbl.Init(nil)
}

func TestDataTable_RenderIdempotent(t *testing.T) {
	tbl := &DataTable{Headers: []string{"A"}, Rows: [][]string{{"1"}}}
	tbl.Init(nil)

	first := tbl.Render().String()
	second := tbl.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
