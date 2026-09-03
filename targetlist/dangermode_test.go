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

// The tick is the always-painted glyph this whole change fixes: it must not
// exist in the stylesheet without a state gate, and the gate must open for
// both selection colours (Selected blue, Invalid red).
func TestTargetList_CheckIconHiddenUntilMarked(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	i := strings.Index(cssStr, ".targetlist__check-icon {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__check-icon")
	}
	body := cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "display: none") {
		t.Errorf("the tick must be hidden by default, block:\n%s", body)
	}
	for _, sel := range []string{
		`[data-selected="true"] .targetlist__check-icon`,
		`[data-invalid="true"] .targetlist__check-icon`,
	} {
		if !strings.Contains(cssStr, sel) {
			t.Errorf("expected the tick to be revealed by %s", sel)
		}
	}
}

// The check rides the row's top-end corner out of the flow: the label never
// shifts when selection mode opens.
func TestTargetList_CheckRidesTopEndCorner(t *testing.T) {
	cssStr := (&TargetList{}).RenderCSS().String()

	i := strings.Index(cssStr, ".targetlist__check {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__check")
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
