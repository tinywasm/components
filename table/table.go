package table

import (
	"github.com/tinywasm/dom"
)

type Table struct {
	*dom.Element
	Headers []string
	Rows    [][]string
}

func (t *Table) Render() *dom.Element {
	if t.Element == nil {
		t.Element = &dom.Element{}
	}

	table := dom.Table().
		Class("table")

	// Header
	thead := dom.Thead()
	tr := dom.Tr()
	for _, header := range t.Headers {
		tr.Add(dom.Th().Text(header))
	}
	thead.Add(tr)
	table.Add(thead)

	// Body
	tbody := dom.Tbody()
	for _, row := range t.Rows {
		tr := dom.Tr()
		for _, cell := range row {
			tr.Add(dom.Td().Text(cell))
		}
		tbody.Add(tr)
	}
	table.Add(tbody)

	return table
}
