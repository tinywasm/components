package datatable

import (
	. "webtyp.com/dom"
	. "webtyp.com/fmt"
	. "webtyp.com/html"
	"webtyp.com/widget"
)

// NameDataTable is the widget name.
const NameDataTable = widget.Name("datatable")

const (
	PartHeader = widget.Part("header")
	PartRow    = widget.Part("row")
)

var (
	clsTable = NameDataTable.Root()
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

func (t *DataTable) WidgetName() widget.Name { return NameDataTable }
func (t *DataTable) WidgetKind() widget.Kind { return widget.Grid }

func (t *DataTable) Init(_ Ctx) {
	t.rows = NewNodes(t.buildRows(t.Rows)...)
}

func (t *DataTable) Render() *Element {
	table := Table().Set(clsTable.AsAttr()).Attr("role", "grid")

	thead := Thead()
	tr := Tr().Attr("role", "row")
	for _, header := range t.Headers {
		tr.Child(Th().Set(NameDataTable.Class(PartHeader).AsAttr()).Attr("role", "columnheader").Text(header))
	}
	thead.Child(tr)
	table.Child(thead)

	// Seeded from the DATA, not from t.rows. The signal owns the live children
	// — BindChildren replaces them wholesale on every update — and its elements
	// already belong to whichever tbody the runtime last painted. Handing them
	// to a second Render() would give one element two parents, so the markup a
	// re-render produces is built fresh. SetRows keeps t.Rows in step, so the
	// seed and the signal always describe the same table.
	tbody := Tbody().BindChildren(t.rows)
	for _, row := range t.buildRows(t.Rows) {
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
		tr := Tr().Set(NameDataTable.Class(PartRow).AsAttr()).Key(Sprint(i)).Attr("role", "row")
		for _, cell := range row {
			tr.Child(Td().Attr("role", "gridcell").Text(cell))
		}
		out = append(out, tr)
	}
	return out
}
