//go:build !wasm

package targetdate_test

import (
	"strings"
	"testing"

	"github.com/tinywasm/components/targetdate"
)

func TestRowsCarryNoOptionsMenu(t *testing.T) {
	td := &targetdate.TargetDate{}
	td.Init(nil)
	td.SetItems([]targetdate.Item{{ID: "1", Label: "Row 1"}})

	html := td.Render().String()

	for _, unwanted := range []string{"targetdate__button", "targetdate__options", "targetdate__item-danger", "Eliminar"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("rendered row contains unwanted options menu artifact %q\nhtml: %s", unwanted, html)
		}
	}
}

func TestCheckIsInTheMarkupWhenModeIsOff(t *testing.T) {
	td := &targetdate.TargetDate{}
	td.Init(nil)
	td.SetItems([]targetdate.Item{{ID: "1", Label: "Row 1"}})

	td.SetSelectMode(false)
	html := td.Render().String()

	if !strings.Contains(html, "targetdate") {
		t.Errorf("container must exist in markup\nhtml: %s", html)
	}
}

func TestSheetValidates(t *testing.T) {
	td := &targetdate.TargetDate{}
	td.Init(nil)
	if errs := td.RenderCSS().String(); errs == "" {
		t.Error("RenderCSS must return stylesheet string")
	}
}
