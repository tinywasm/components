package statgrid

import (
	. "webtyp.com/dom"
	. "webtyp.com/html"
	"webtyp.com/widget"
)

// NameStatGrid is the widget name for statgrid.
const NameStatGrid = widget.Name("statgrid")

const (
	PartItem  = widget.Part("item")
	PartValue = widget.Part("value")
	PartLabel = widget.Part("label")
)

var (
	clsStatGrid  = NameStatGrid.Root()
	clsStatItem  = NameStatGrid.Class(PartItem)
	clsStatValue = NameStatGrid.Class(PartValue)
	clsStatLabel = NameStatGrid.Class(PartLabel)
)

// StatItem represents a single statistic pair (value and label).
type StatItem struct {
	Value string
	Label string
}

// StatGrid displays a responsive grid of key metrics/statistics.
type StatGrid struct {
	Element
	Items []StatItem
}

func (s *StatGrid) WidgetName() widget.Name { return NameStatGrid }
func (s *StatGrid) WidgetKind() widget.Kind { return widget.Region }

func (s *StatGrid) Render() *Element {
	grid := Div().Set(clsStatGrid.AsAttr())

	for _, item := range s.Items {
		card := Div().Set(clsStatItem.AsAttr()).Child(
			Div().Set(clsStatValue.AsAttr()).Text(item.Value),
			Div().Set(clsStatLabel.AsAttr()).Text(item.Label),
		)
		grid.Child(card)
	}

	return grid
}
