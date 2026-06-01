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
	table := Table().Add(clsTable.AsAttr())

	// Header
	thead := Thead()
	tr := Tr()
	for _, header := range t.Headers {
		tr.Add(Th().Text(header))
	}
	thead.Add(tr)
	table.Add(thead)

	// Body
	tbody := Tbody()
	for _, row := range t.Rows {
		tr := Tr()
		for _, cell := range row {
			tr.Add(Td().Text(cell))
		}
		tbody.Add(tr)
	}
	table.Add(tbody)

	return table
}
