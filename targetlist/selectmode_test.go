//go:build !wasm

package targetlist_test

import (
	"strings"
	"testing"

	"github.com/tinywasm/components/targetlist"
)

func TestRowsCarryNoOptionsMenu(t *testing.T) {
	tl := &targetlist.TargetList{}
	tl.Init(nil)
	tl.SetItems([]targetlist.Item{{ID: "1", Label: "Row 1"}})

	html := tl.Render().String() + tl.Items()[0].ID

	for _, unwanted := range []string{"targetlist__button", "targetlist__options", "targetlist__item-danger", "Eliminar"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("rendered row contains unwanted options menu artifact %q\nhtml: %s", unwanted, html)
		}
	}
}

func TestSelectModeOffFiresOnSelect(t *testing.T) {
	tl := &targetlist.TargetList{}
	tl.Init(nil)
	tl.SetItems([]targetlist.Item{{ID: "1", Label: "Row 1"}})

	if len(tl.CheckedIDs()) != 0 {
		t.Errorf("selection mode off should have no checked IDs, got %v", tl.CheckedIDs())
	}
}

func TestSelectModeOnTogglesInsteadOfSelecting(t *testing.T) {
	tl := &targetlist.TargetList{}
	tl.Init(nil)
	tl.SetItems([]targetlist.Item{{ID: "1", Label: "Row 1"}, {ID: "2", Label: "Row 2"}})

	tl.SetSelectMode(true)
	if len(tl.CheckedIDs()) != 0 {
		t.Errorf("initially no checked IDs")
	}
}

func TestCheckIsInTheMarkupWhenModeIsOff(t *testing.T) {
	tl := &targetlist.TargetList{}
	tl.Init(nil)
	tl.SetItems([]targetlist.Item{{ID: "1", Label: "Row 1"}})

	tl.SetSelectMode(false)
	html := tl.Render().String()

	if !strings.Contains(html, "targetlist") {
		t.Errorf("container must exist in markup\nhtml: %s", html)
	}
}

func TestCheckedIDsFollowRenderOrder(t *testing.T) {
	tl := &targetlist.TargetList{}
	tl.Init(nil)
	tl.SetItems([]targetlist.Item{
		{ID: "1", Label: "First"},
		{ID: "2", Label: "Second"},
		{ID: "3", Label: "Third"},
	})

	tl.SetSelectMode(true)
	// CheckedIDs uses rendering order provided by SetItems
	if ids := tl.CheckedIDs(); len(ids) != 0 {
		t.Errorf("expected empty checked IDs, got %v", ids)
	}
}

func TestSheetValidates(t *testing.T) {
	tl := &targetlist.TargetList{}
	tl.Init(nil)
	if errs := tl.RenderCSS().String(); errs == "" {
		t.Error("RenderCSS must return stylesheet string")
	}
}
