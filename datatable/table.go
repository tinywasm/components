package datatable

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

var (
	clsTable Class = "table"
)

type DataTable struct {
	Element
	Headers []string
	Rows    [][]string
}

func (t *DataTable) Render() *Element {
	table := Table().Set(clsTable.AsAttr())

	// Header
	thead := Thead()
	tr := Tr()
	for _, header := range t.Headers {
		tr.Child(Th().Text(header))
	}
	thead.Child(tr)
	table.Child(thead)

	// Body
	tbody := Tbody()
	for _, row := range t.Rows {
		tr := Tr()
		for _, cell := range row {
			tr.Child(Td().Text(cell))
		}
		tbody.Child(tr)
	}
	table.Child(tbody)

	return table
}
