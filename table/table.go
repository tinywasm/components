package table

import (
	"github.com/tinywasm/components"
	"github.com/tinywasm/dom"
)

type Table struct {
	components.BaseComponent
	Headers []string
	Rows    [][]string
}

func (t *Table) Render() dom.Node {
	table := dom.Tag("table").
		ID(t.ID()).
		Class("table")

	// Header
	thead := dom.Tag("thead")
	tr := dom.Tag("tr")
	for _, header := range t.Headers {
		tr.Append(dom.Tag("th").Text(header))
	}
	thead.Append(tr)
	table.Append(thead)

	// Body
	tbody := dom.Tag("tbody")
	for _, row := range t.Rows {
		tr := dom.Tag("tr")
		for _, cell := range row {
			tr.Append(dom.Tag("td").Text(cell))
		}
		tbody.Append(tr)
	}
	table.Append(tbody)

	return table.ToNode()
}
