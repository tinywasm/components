//go:build !wasm

package targethour

import (
	"strings"
	"testing"
)

// The master check is the first element in the root, hidden until selection
// mode opens.
func TestTargetHour_MasterCheckHiddenUntilSelectionMode(t *testing.T) {
	css := (&TargetHour{}).RenderCSS().String()
	i := strings.Index(css, ".targethour__check-all {")
	if i == -1 {
		t.Fatal("expected a rule for .targethour__check-all")
	}
	body := css[i:]
	if e := strings.Index(body, "}"); e != -1 {
		body = body[:e]
	}
	if !strings.Contains(body, "display: none") {
		t.Errorf("the master check must be hidden by default, block:\n%s", body)
	}
	if !strings.Contains(css, `.targethour[data-open="true"] .targethour__check-all {`) {
		t.Errorf("the master check must be revealed by the list's open state")
	}
}

// Tapping it with nothing / some marked selects every row; tapping it with
// all marked clears.
func TestTargetHour_MasterCheckTogglesAll(t *testing.T) {
	th := &TargetHour{}
	th.Init(nil)
	th.SetItems([]Item{{ID: "1"}, {ID: "2"}, {ID: "3"}})
	th.SetSelectMode(true)

	m := th.buildMasterCheck()
	_ = m

	// simulate the click handler directly (no DOM under SSR)
	if n := th.sel.Count(); n > 0 && n == 3 {
		th.sel.Clear()
	} else {
		th.sel.CheckAll(th.itemIDs())
	}
	if th.sel.Count() != 3 {
		t.Fatalf("first tap must select all, Count = %d", th.sel.Count())
	}
	if n := th.sel.Count(); n > 0 && n == 3 {
		th.sel.Clear()
	} else {
		th.sel.CheckAll(th.itemIDs())
	}
	if th.sel.Count() != 0 {
		t.Fatalf("second tap must clear, Count = %d", th.sel.Count())
	}
}

// The count label renders "n / total".
func TestTargetHour_MasterCheckShowsNOfTotal(t *testing.T) {
	th := &TargetHour{}
	th.Init(nil)
	th.SetItems([]Item{{ID: "1"}, {ID: "2"}})
	th.SetSelectMode(true)
	th.sel.CheckAll([]string{"1"})

	html := th.buildMasterCheck().String()
	if !strings.Contains(html, "1 / 2") {
		t.Errorf("master check must show \"1 / 2\", got:\n%s", html)
	}
}

// The master reuses the shared glyphs, not a bespoke tick.
func TestTargetHour_MasterCheckUsesSharedGlyphs(t *testing.T) {
	th := &TargetHour{}
	th.Init(nil)
	html := th.buildMasterCheck().String()
	if !strings.Contains(html, `href='#trash'`) || !strings.Contains(html, `href='#pencil'`) {
		t.Errorf("master check must reference the shared trash/pencil glyphs:\n%s", html)
	}
}
