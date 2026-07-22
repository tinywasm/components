package targetlist

import (
	"strings"
	"testing"
)

func TestTargetList_RowHasLabelBadgeAndMenu(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)

	html := tl.buildRow(Item{ID: "7", Label: "Alpha", Description: "192.168.0.7"}).String()

	for _, want := range []string{"tl-row", "Alpha", "tl-badge", "192.168.0.7", "tl-menu", "Editar", "Eliminar"} {
		if !strings.Contains(html, want) {
			t.Errorf("buildRow output missing %q\ngot: %s", want, html)
		}
	}
}

func TestTargetList_SetItemsPopulatesRows(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1", Label: "One"}, {ID: "2", Label: "Two"}})

	if got := len(tl.rows.Get()); got != 2 {
		t.Fatalf("expected 2 rows, got %d", got)
	}
}
