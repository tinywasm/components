//go:build !wasm

package targetlist

import (
	"regexp"
	"strings"
	"testing"
)

var (
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
)

func TestTargetList_RowHasLabelBadgeAndMenu(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)

	html := tl.buildRow(Item{ID: "7", Label: "Alpha", Description: "192.168.0.7"}).String()

	for _, want := range []string{"targetlist__row", "Alpha", "targetlist__badge", "192.168.0.7", "targetlist__menu", "Editar", "Eliminar"} {
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

func TestTargetList_MenuOpenStateBackdrop(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1", Label: "One"}})

	// Initially, with no open menus, backdrop shouldn't have data-open="true"
	htmlInit := tl.Render().String()
	t.Logf("htmlInit: %s", htmlInit)
	if strings.Contains(htmlInit, `data-open='true'`) {
		t.Error("expected backdrop NOT to have data-open='true' initially")
	}

	// Mocking menu open
	tl.menuOpen.Set(true)
	htmlOpen := tl.Render().String()
	t.Logf("htmlOpen: %s", htmlOpen)
	if !strings.Contains(htmlOpen, `data-open='true'`) {
		t.Error("expected backdrop to have data-open='true' when a menu is open")
	}

	// Verify closeAllMenus clears it
	tl.closeAllMenus()
	if tl.menuOpen.Get() {
		t.Error("expected menuOpen to be false after closeAllMenus")
	}
}

func TestTargetList_CSSDoesNotContainHas(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	css := tl.RenderCSS().String()

	if strings.Contains(css, ":has(") {
		t.Error("expected CSS not to contain forbidden :has( selector")
	}
	if !strings.Contains(css, "display: none") {
		t.Error("expected CSS to have display: none under normal condition for backdrop")
	}
	if !strings.Contains(css, "[data-open=\"true\"]") {
		t.Error("expected CSS to contain selector matching [data-open=\"true\"]")
	}
	if !strings.Contains(css, "display: block") {
		t.Error("expected CSS to have display: block under open condition")
	}
}

func TestPairMarkupAndStylesheet(t *testing.T) {
	extractCSSClasses := func(css string) map[string]bool {
		classes := make(map[string]bool)
		matches := classRegex.FindAllStringSubmatch(css, -1)
		for _, m := range matches {
			if len(m) > 1 {
				classes[m[1]] = true
			}
		}
		return classes
	}

	extractHTMLClasses := func(html string) map[string]bool {
		classes := make(map[string]bool)
		matches1 := htmlClassRegex.FindAllStringSubmatch(html, -1)
		for _, m := range matches1 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		matches2 := htmlClassRegex2.FindAllStringSubmatch(html, -1)
		for _, m := range matches2 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		return classes
	}

	filterClasses := func(classes map[string]bool, prefix string) map[string]bool {
		filtered := make(map[string]bool)
		for cls := range classes {
			if strings.HasPrefix(cls, prefix) {
				filtered[cls] = true
			}
		}
		return filtered
	}

	tl := &TargetList{}
	tl.Init(nil)
	html := tl.Render().String() + tl.buildRow(Item{ID: "1", Label: "A", Description: "B"}).String()
	css := tl.RenderCSS().String()

	htmlClasses := filterClasses(extractHTMLClasses(html), "targetlist")
	cssClasses := filterClasses(extractCSSClasses(css), "targetlist")

	for cls := range cssClasses {
		if !htmlClasses[cls] {
			t.Errorf("TargetList CSS class %q does not exist in rendered HTML", cls)
		}
	}
	for cls := range htmlClasses {
		if !cssClasses[cls] {
			t.Errorf("TargetList HTML class %q is unstyled in CSS", cls)
		}
	}
}
