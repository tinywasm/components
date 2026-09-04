//go:build !wasm

package targetdate

import (
	"strings"
	"testing"
)

// Internal, so it can reach buildRow: Render() binds its children and SSR does
// not serialize a children binding, so this is the only place a row's markup
// can actually be asserted. targetlist carries the same test — the two
// components must stay interchangeable for crudview, and a check that exists
// in one and not the other is exactly how they drift apart.
func TestTargetDate_RowHasLeadLabelAndCheck(t *testing.T) {
	td := &TargetDate{}
	td.Init(nil)

	html := td.buildRow(Item{
		ID: "7", Label: "Fractura", Description: "dr. Tony Stark",
		LeadTop: "Jue", LeadMain: "15", LeadBottom: "Ene 26",
	}).String()

	for _, want := range []string{"targetdate__row", "Fractura", "targetdate__sel-check", "15"} {
		if !strings.Contains(html, want) {
			t.Errorf("buildRow output missing %q\ngot: %s", want, html)
		}
	}
	for _, unwanted := range []string{"targetdate__button", "targetdate__options", "targetdate__item-danger", "Eliminar"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("buildRow must not render per-row menu artifact %q\ngot: %s", unwanted, html)
		}
	}
}
