package datatable

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
)

var (
	clsTable Class = "table"
)

// DataTable renders a table whose body rows can be provided up front
// (Rows) and/or updated later via SetRows (e.g. after an async fetch
// resolves). Header columns are static — only the body is reactive.
type DataTable struct {
	Element            // value embed — NEVER pointer (TinyGo heap constraint)
	Headers []string   // column titles, fixed at construction
	Rows    [][]string // initial rows, rendered on first paint; empty is fine

	rows *SignalNodes // internal — drives the reactive tbody
}

func (t *DataTable) Init(_ Ctx) {
	t.rows = NewNodes(t.buildRows(t.Rows)...)
}

func (t *DataTable) Render() *Element {
	table := Table().Set(clsTable.AsAttr())

	thead := Thead()
	tr := Tr()
	for _, header := range t.Headers {
		tr.Child(Th().Text(header))
	}
	thead.Child(tr)
	table.Child(thead)

	tbody := Tbody().BindChildren(t.rows)
	for _, row := range t.rows.Get() {
		tbody.Child(row)
	}
	table.Child(tbody)

	return table
}

// SetRows replaces the table body — safe to call after Init/Render,
// e.g. once data from an async source (fetch, MCP call) arrives.
func (t *DataTable) SetRows(rows [][]string) {
	t.Rows = rows
	t.rows.Set(t.buildRows(rows))
}

func (t *DataTable) buildRows(rows [][]string) []*Element {
	out := make([]*Element, 0, len(rows))
	for i, row := range rows {
		tr := Tr().Key(Sprint(i))
		for _, cell := range row {
			tr.Child(Td().Text(cell))
		}
		out = append(out, tr)
	}
	return out
}
