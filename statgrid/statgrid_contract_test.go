package statgrid

import (
	"testing"

	"webtyp.com/dom"
)

var _ dom.Component = (*StatGrid)(nil)

func TestStatGrid_RenderIdempotent(t *testing.T) {
	s := &StatGrid{
		Items: []StatItem{
			{Value: "80+", Label: "Años de historia"},
			{Value: "15k", Label: "Atenciones al año"},
		},
	}

	first := s.Render().String()
	second := s.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
