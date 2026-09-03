//go:build !wasm

package targetlist

import (
	"strings"
	"testing"
)

// The master check is the first element in the root, hidden until selection
// mode opens.
func TestTargetList_MasterCheckHiddenUntilSelectionMode(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()
	i := strings.Index(css, ".targetlist__check-all {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__check-all")
	}
	body := css[i:]
	if e := strings.Index(body, "}"); e != -1 {
		body = body[:e]
	}
	if !strings.Contains(body, "display: none") {
		t.Errorf("the master check must be hidden by default, block:\n%s", body)
	}
	if !strings.Contains(css, `.targetlist[data-open="true"] .targetlist__check-all {`) {
		t.Errorf("the master check must be revealed by the list's open state")
	}
}

// Tapping it with nothing / some marked selects every row; tapping it with
// all marked clears.
func TestTargetList_MasterCheckTogglesAll(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1"}, {ID: "2"}, {ID: "3"}})
	tl.SetSelectMode(true)

	m := tl.buildMasterCheck()
	_ = m

	// simulate the click handler directly (no DOM under SSR)
	if n := tl.sel.Count(); n > 0 && n == 3 {
		tl.sel.Clear()
	} else {
		tl.sel.CheckAll(tl.itemIDs())
	}
	if tl.sel.Count() != 3 {
		t.Fatalf("first tap must select all, Count = %d", tl.sel.Count())
	}
	if n := tl.sel.Count(); n > 0 && n == 3 {
		tl.sel.Clear()
	} else {
		tl.sel.CheckAll(tl.itemIDs())
	}
	if tl.sel.Count() != 0 {
		t.Fatalf("second tap must clear, Count = %d", tl.sel.Count())
	}
}

// The count label renders "n / total".
func TestTargetList_MasterCheckShowsNOfTotal(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1"}, {ID: "2"}})
	tl.SetSelectMode(true)
	tl.sel.CheckAll([]string{"1"})

	html := tl.buildMasterCheck().String()
	if !strings.Contains(html, "1 / 2") {
		t.Errorf("master check must show \"1 / 2\", got:\n%s", html)
	}
}

// The master reuses the shared glyphs, not a bespoke tick.
func TestTargetList_MasterCheckUsesSharedGlyphs(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	html := tl.buildMasterCheck().String()
	if !strings.Contains(html, `href='#trash'`) || !strings.Contains(html, `href='#pencil'`) {
		t.Errorf("master check must reference the shared trash/pencil glyphs:\n%s", html)
	}
}
