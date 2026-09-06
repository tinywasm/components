//go:build !wasm

package targetlist

import (
	"strings"
	"testing"
)

// A checked row under the armed danger tone is Invalid (red) INSTEAD OF
// Selected (blue): one row, one fill, so the two stylesheet rules cannot
// race in the cascade. Without the tone it stays Selected as always.
func TestTargetList_DangerMarksInvalidNotSelected(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	row := tl.buildRow(Item{ID: "1", Label: "Row 1"})

	tl.SetSelectMode(true)
	tl.SetDanger(true)
	tl.sel.Toggle("1")

	html := row.String()
	if !strings.Contains(html, `data-invalid='true'`) {
		t.Errorf("a checked row with the danger tone must carry data-invalid\nhtml: %s", html)
	}
	if strings.Contains(html, "data-selected") {
		t.Errorf("a danger-marked row must not also carry data-selected\nhtml: %s", html)
	}

	tl.SetDanger(false)
	html = row.String()
	if !strings.Contains(html, `data-selected='true'`) {
		t.Errorf("without the tone a checked row stays selected\nhtml: %s", html)
	}
	if strings.Contains(html, "data-invalid") {
		t.Errorf("without the tone no row may carry data-invalid\nhtml: %s", html)
	}
}

// THE regression: a row that is merely the loaded record in NORMAL mode
// (Selected on the <li>, no selection mode) must NOT put any state on the
// check box — otherwise the glyph-reveal rule (.check[data-selected]
// .check-pencil) fires and a pencil appears on a plain highlight. The check
// box's state is written only in selection mode.
func TestTargetList_NormalHighlightLeavesCheckStateless(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.Selected.Set("1") // loaded in the form — normal-mode highlight, no select mode
	row := tl.buildRow(Item{ID: "1", Label: "Row 1"})

	html := row.String()
	check := html[strings.Index(html, `class='targetlist__sel-check'`):]
	check = check[:strings.Index(check, "</span>")]
	if strings.Contains(check, "data-selected") || strings.Contains(check, "data-invalid") {
		t.Errorf("normal-mode highlight must leave the check box stateless\ncheck: %s", check)
	}
	// The <li> itself still carries the highlight.
	if !strings.Contains(html, `data-selected='true'`) {
		t.Errorf("the row <li> must still carry the normal-mode highlight\nhtml: %s", html)
	}
}

// Which glyph shows hangs off the CHECK BOX's own state, never the row's:
// .check[data-invalid] reveals the trash, .check[data-selected] the pencil.
// Both glyphs hidden by default.
func TestTargetList_GlyphRevealHangsOffTheCheckBox(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	for _, part := range []string{".targetlist__sel-check-trash", ".targetlist__sel-check-pencil"} {
		i := strings.Index(cssStr, part+" {")
		if i == -1 {
			t.Fatalf("expected a rule for %s", part)
		}
		body := cssStr[i:]
		if end := strings.Index(body, "}"); end != -1 {
			body = body[:end]
		}
		if !strings.Contains(body, "display: none") {
			t.Errorf("%s must be hidden by default, block:\n%s", part, body)
		}
	}

	for _, sel := range []string{
		`.targetlist__sel-check[data-invalid="true"] .targetlist__sel-check-trash`,
		`.targetlist__sel-check[data-selected="true"] .targetlist__sel-check-pencil`,
	} {
		if !strings.Contains(cssStr, sel) {
			t.Errorf("expected the glyph to be revealed by %s", sel)
		}
	}
	// It must NOT be revealed from the row: that state is also the normal-mode
	// highlight.
	for _, sel := range []string{
		`.targetlist__row[data-selected="true"] .targetlist__sel-check-pencil`,
		`.targetlist__row[data-invalid="true"] .targetlist__sel-check-trash`,
	} {
		if strings.Contains(cssStr, sel) {
			t.Errorf("the glyph reveal must not hang off the row: %s", sel)
		}
	}
}

// The box fills white-on-colour for whichever action: AccentInverse (amber,
// white on-primary) for edit, Danger (red, white on-danger) for delete. An
// Accent/Inset fill would paint the glyph near-black through currentColor —
// the exact defect this fixes.
func TestTargetList_CheckBoxCarriesWhiteText(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	for _, tc := range []struct{ sel, want string }{
		{`.targetlist__sel-check[data-selected="true"]`, "--color-on-primary"},
		{`.targetlist__sel-check[data-invalid="true"]`, "--color-on-danger"},
	} {
		i := strings.Index(cssStr, tc.sel+" {")
		if i == -1 {
			t.Fatalf("expected a rule for %s", tc.sel)
		}
		body := cssStr[i:]
		if end := strings.Index(body, "}"); end != -1 {
			body = body[:end]
		}
		if !strings.Contains(body, tc.want) {
			t.Errorf("expected %s in the %s rule (glyph must be white), block:\n%s", tc.want, tc.sel, body)
		}
		if strings.Contains(body, "--color-on-accent") {
			t.Errorf("the check rule must not use on-accent (near-black glyph), block:\n%s", body)
		}
	}
}

// The row markup references the shared webtyp/icons glyphs, not a fixed tick.
func TestTargetList_RenderUsesSharedIconGlyphs(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	row := tl.buildRow(Item{ID: "1", Label: "Row 1"})

	html := row.String()
	if !strings.Contains(html, `href='#trash'`) || !strings.Contains(html, `href='#pencil'`) {
		t.Errorf("the row must reference the shared trash/pencil glyphs\nhtml: %s", html)
	}
	if strings.Contains(html, "tl-check") || strings.Contains(html, "crud-delete") {
		t.Errorf("the old fixed/crudicons glyph ids must be gone\nhtml: %s", html)
	}
}

// The box is a centred flex square, revealed only inside the list root's Open
// state (selection mode). In normal mode the root has no data-open, so the
// square does not render at all — not an empty one.
func TestTargetList_CheckHiddenUntilSelectionMode(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	// Anchored at line start: the KeepSize primitives group also mentions
	// this class (as its last selector), so a bare substring search would
	// land on a rule that has flex-shrink/grow but no display authority.
	i := strings.Index(cssStr, "\n.targetlist__sel-check {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__sel-check")
	}
	body := cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "display: none") {
		t.Errorf("the check box must be hidden by default (normal mode), block:\n%s", body)
	}

	i = strings.Index(cssStr, `.targetlist[data-open="true"] .targetlist__sel-check {`)
	if i == -1 {
		t.Fatal("expected the check box to be revealed by the list's open state")
	}
	body = cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "display: flex") || !strings.Contains(body, "justify-content: center") {
		t.Errorf("the revealed box must be a centred flex box (glyph dead centre), block:\n%s", body)
	}

	// Each glyph carries its own size — an svg sized only by the box overflows it.
	for _, part := range []string{".targetlist__sel-check-trash", ".targetlist__sel-check-pencil"} {
		i := strings.Index(cssStr, part+" {")
		if i == -1 {
			t.Fatalf("expected a rule for %s", part)
		}
		body := cssStr[i:]
		if end := strings.Index(body, "}"); end != -1 {
			body = body[:end]
		}
		if !strings.Contains(body, "width:") || !strings.Contains(body, "height:") {
			t.Errorf("%s must carry its own IconBox size, block:\n%s", part, body)
		}
	}
}

// The check rides the row's top-end corner out of the flow: the label never
// shifts when selection mode opens.
func TestTargetList_CheckRidesTopEndCorner(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	i := strings.Index(cssStr, "\n.targetlist__sel-check {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__sel-check")
	}
	body := cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "position: absolute") {
		t.Errorf("the check must be out of flow (position: absolute), block:\n%s", body)
	}
}

// The danger row is a wash — red like the hover, not a solid fill.
func TestTargetList_DangerRowUsesWash(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	i := strings.Index(cssStr, `.targetlist__row[data-invalid="true"] {`)
	if i == -1 {
		t.Fatal("expected a rule for the invalid row state")
	}
	body := cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "--color-danger-wash") {
		t.Errorf("a danger-marked row must use the DangerWash surface, block:\n%s", body)
	}
}
