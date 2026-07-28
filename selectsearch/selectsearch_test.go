//go:build !wasm

package selectsearch

import (
	"regexp"
	"strings"
	"testing"

	"github.com/tinywasm/fmt"
)

var (
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
)

func TestSelectSearch_Render(t *testing.T) {
	c := &SelectSearch{
		Placeholder: "Choose category",
		Options: []SsOption{
			{ID: "1", Label: "Automobiles", Description: "auto"},
			{ID: "2", Label: "Film & Animation", Description: "anime"},
		},
	}
	c.Init(nil)

	html := c.Render().String()

	if !fmt.Contains(html, "selectsearch") {
		t.Error("expected selectsearch class")
	}
	if !fmt.Contains(html, "selectsearch__toggle") {
		t.Error("expected selectsearch__toggle checkbox")
	}
	if !fmt.Contains(html, "Choose category") {
		t.Error("expected placeholder text")
	}
	// Note: options are in the dropdown, which is wrapped in Show(c.isOpen)
	// Since isOpen is false, options won't be in the initial static HTML
}

func TestSelectSearch_SelectedValue(t *testing.T) {
	c := &SelectSearch{
		Placeholder: "Choose category",
	}
	c.Init(nil)
	c.selectedLabel.Set("Automobiles")

	html := c.Render().String()
	if !fmt.Contains(html, "Automobiles") {
		t.Error("expected selected label")
	}
}

func TestSelectSearch_OpenState_RendersChecked(t *testing.T) {
	c := &SelectSearch{
		Options: []SsOption{
			{ID: "1", Label: "Apple"},
		},
	}
	c.Init(nil)
	c.isOpen.Set(true)

	el := c.Render()
	html := el.String()
	if !fmt.Contains(html, "checked") {
		t.Error("expected 'checked' attribute on toggle when isOpen=true")
	}
}

func TestSelectSearch_Filtering(t *testing.T) {
	c := &SelectSearch{
		Options: []SsOption{
			{ID: "1", Label: "Automobiles"},
			{ID: "2", Label: "Film & Animation"},
		},
	}
	c.Init(nil)
	c.isOpen.Set(true)
	c.query.Set("Film")
	// query.Set doesn't trigger OnChange in standard tests unless manually called or using gotest/WASM
	c.rows.Set(c.buildRows("Film"))

	el := c.Render()
	html := el.String()
	if fmt.Contains(html, "Automobiles") {
		t.Error("expected Automobiles to be filtered out")
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

	ss := &SelectSearch{}
	ss.Init(nil)
	ss.isOpen.Set(true)
	ss.SetOptions([]SsOption{{ID: "1", Label: "A", Description: "B"}})
	html := ss.Render().String()
	// Render option row too
	html += ss.buildRows("")[0].String()

	css := ss.RenderCSS().String()

	htmlClasses := filterClasses(extractHTMLClasses(html), "selectsearch")
	cssClasses := filterClasses(extractCSSClasses(css), "selectsearch")

	for cls := range cssClasses {
		if !htmlClasses[cls] {
			t.Errorf("SelectSearch CSS class %q does not exist in rendered HTML", cls)
		}
	}
	for cls := range htmlClasses {
		if !cssClasses[cls] {
			t.Errorf("SelectSearch HTML class %q is unstyled in CSS", cls)
		}
	}
}
