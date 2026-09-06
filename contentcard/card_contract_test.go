package contentcard

import (
	"testing"

	"webtyp.com/dom"
)

var _ dom.Component = (*ContentCard)(nil)

func TestContentCard_RenderIdempotent(t *testing.T) {
	c := &ContentCard{}

	first := c.Render().String()
	second := c.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
