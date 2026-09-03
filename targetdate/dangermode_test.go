//go:build !wasm

package targetdate

import (
	"strings"
	"testing"
)

// Mirror of targetlist's danger contract: the two lists stay
// interchangeable for crudview, and a state that exists in one and not the
// other is exactly how they drift apart.
func TestTargetDate_DangerMarksInvalidNotSelected(t *testing.T) {
	td := &TargetDate{}
	td.Init(nil)
	row := td.buildRow(Item{ID: "1", Label: "Row 1"})

	td.SetSelectMode(true)
	td.SetDanger(true)
	td.sel.Toggle("1")

	html := row.String()
	if !strings.Contains(html, `data-invalid='true'`) {
		t.Errorf("a checked row with the danger tone must carry data-invalid\nhtml: %s", html)
	}
	if strings.Contains(html, "data-selected") {
		t.Errorf("a danger-marked row must not also carry data-selected\nhtml: %s", html)
	}

	td.SetDanger(false)
	html = row.String()
	if !strings.Contains(html, `data-selected='true'`) {
		t.Errorf("without the tone a checked row stays selected\nhtml: %s", html)
	}
	if strings.Contains(html, "data-invalid") {
		t.Errorf("without the tone no row may carry data-invalid\nhtml: %s", html)
	}
}

func TestTargetDate_CheckRidesTopEndCorner(t *testing.T) {
	cssStr := (&TargetDate{}).RenderCSS().String()

	i := strings.Index(cssStr, ".targetdate__check {")
	if i == -1 {
		t.Fatal("expected a rule for .targetdate__check")
	}
	body := cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "position: absolute") {
		t.Errorf("the check must be out of flow (position: absolute), block:\n%s", body)
	}
}

func TestTargetDate_DangerRowUsesWash(t *testing.T) {
	cssStr := (&TargetDate{}).RenderCSS().String()

	i := strings.Index(cssStr, `.targetdate__row[data-invalid="true"] {`)
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
