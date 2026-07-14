package actionbutton

import (
	"testing"

	"github.com/tinywasm/dom"
)

var _ dom.Component = (*ActionButton)(nil)

func TestActionButton_RenderIdempotent(t *testing.T) {
	b := &ActionButton{Text: "Save", Variant: "primary"}

	first := b.Render().String()
	second := b.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
