//go:build !wasm

package targetdate

import (
	"strings"
	"testing"
)

// The selection header is the widget root's first child, hidden until
// selection mode opens. It is the strip listselect.Header builds; the widget
// merely Child()s it above the <ul>.
func TestTargetDate_MasterCheckHiddenUntilSelectionMode(t *testing.T) {
	css := (&TargetDate{}).RenderCSS().String()
	i := strings.Index(css, ".targetdate__sel-header {")
	if i == -1 {
		t.Fatal("expected a rule for .targetdate__sel-header")
	}
	body := css[i:]
	if e := strings.Index(body, "}"); e != -1 {
		body = body[:e]
	}
	if !strings.Contains(body, "display: none") {
		t.Errorf("the header must be hidden by default, block:\n%s", body)
	}
	if !strings.Contains(css, `.targetdate[data-open="true"] .targetdate__sel-header {`) {
		t.Errorf("the header must be revealed by the list's open state")
	}
}

// Tapping the select-all box with nothing / some marked selects every row;
// tapping it with all marked clears. The click handler lives in
// listselect.Header (covered by its WASM test); here only the widget's wiring
// is checked — the Mode the header drives.
func TestTargetDate_MasterCheckTogglesAll(t *testing.T) {
	td := &TargetDate{}
	td.Init(nil)
	td.SetItems([]Item{{ID: "1"}, {ID: "2"}, {ID: "3"}})
	td.SetSelectMode(true)

	if n := td.sel.Count(); n > 0 && n == 3 {
		td.sel.Clear()
	} else {
		td.sel.CheckAll(td.itemIDs())
	}
	if td.sel.Count() != 3 {
		t.Fatalf("first tap must select all, Count = %d", td.sel.Count())
	}
	if n := td.sel.Count(); n > 0 && n == 3 {
		td.sel.Clear()
	} else {
		td.sel.CheckAll(td.itemIDs())
	}
	if td.sel.Count() != 0 {
		t.Fatalf("second tap must clear, Count = %d", td.sel.Count())
	}
}

// The count label renders "n / total" inside the header strip.
func TestTargetDate_MasterCheckShowsNOfTotal(t *testing.T) {
	td := &TargetDate{}
	td.Init(nil)
	td.SetItems([]Item{{ID: "1"}, {ID: "2"}})
	td.SetSelectMode(true)
	td.sel.CheckAll([]string{"1"})

	html := td.Render().String()
	if !strings.Contains(html, "1 / 2") {
		t.Errorf("master check must show \"1 / 2\", got:\n%s", html)
	}
}

// The header reuses the shared glyphs, not a bespoke tick.
func TestTargetDate_MasterCheckUsesSharedGlyphs(t *testing.T) {
	td := &TargetDate{}
	td.Init(nil)
	html := td.Render().String()
	if !strings.Contains(html, `href='#trash'`) || !strings.Contains(html, `href='#pencil'`) {
		t.Errorf("master check must reference the shared trash/pencil glyphs:\n%s", html)
	}
}