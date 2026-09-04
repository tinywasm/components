//go:build !wasm

package targethour

import (
	"strings"
	"testing"
)

func TestTargetHour_RowHasHourLabelAndCheck(t *testing.T) {
	th := &TargetHour{}
	th.Init(nil)

	html := th.buildRow(Item{
		ID: "7", Label: "Fractura", Description: "dr. Tony Stark",
		LeadMain: "10:30",
	}).String()

	for _, want := range []string{"targethour__row", "Fractura", "targethour__sel-check", "targethour__hour", "10:30"} {
		if !strings.Contains(html, want) {
			t.Errorf("buildRow output missing %q\ngot: %s", want, html)
		}
	}
	for _, unwanted := range []string{"targethour__button", "targethour__options", "targethour__item-danger", "Eliminar"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("buildRow must not render per-row menu artifact %q\ngot: %s", unwanted, html)
		}
	}
}
