package modaldialog

import (
	"regexp"
	"testing"

	"github.com/tinywasm/dom"
)

var _ dom.Component = (*ModalDialog)(nil)

var idAttr = regexp.MustCompile(`id='[^']*'`)

func TestDialogWidget_InitTwiceSafe(t *testing.T) {
	m := &ModalDialog{Title: "Confirm"}

	m.Init(nil)
	m.Init(nil)
}

// Render is idempotent in shape: each call may assign fresh auto-ids
// (e.g. via Show), so ids are normalized before comparing.
func TestDialogWidget_RenderIdempotent(t *testing.T) {
	m := &ModalDialog{Title: "Confirm"}
	m.Init(nil)

	first := idAttr.ReplaceAllString(m.Render().String(), "id='_'")
	second := idAttr.ReplaceAllString(m.Render().String(), "id='_'")

	if first != second {
		t.Errorf("Render() not idempotent in shape: first=%q second=%q", first, second)
	}
}
