//go:build !wasm

package countbadge

import (
	"testing"

	"webtyp.com/dom"
)

var _ dom.Component = (*CountBadge)(nil)

func TestCountBadge_RenderIdempotent(t *testing.T) {
	b := newBadge("3", true)

	first := b.Render().String()
	second := b.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
