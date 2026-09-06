//go:build !wasm

package datatable

import (
	"regexp"
	"strings"
	"testing"

	. "webtyp.com/fmt"
)

var (
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
)

func TestTable_Render(t *testing.T) {
	tbl := &DataTable{
		Headers: []string{"Name", "Age"},
		Rows: [][]string{
			{"Alice", "30"},
			{"Bob", "25"},
		},
	}
	tbl.Init(nil)

	html := tbl.Render().String()

	if !HasPrefix(html, "<table") {
		t.Error("expected table tag")
	}
	if !Contains(html, "class='datatable'") {
		t.Error("expected datatable class")
	}

	// Check headers
	if !Contains(html, "class='datatable__header' role='columnheader'>Name</th>") {
		t.Error("expected Name header")
	}
	if !Contains(html, "class='datatable__header' role='columnheader'>Age</th>") {
		t.Error("expected Age header")
	}

	// Check rows
	if !Contains(html, "role='gridcell'>Alice</td>") {
		t.Error("expected Alice cell")
	}
	if !Contains(html, "role='gridcell'>30</td>") {
		t.Error("expected 30 cell")
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

	dt := &DataTable{Headers: []string{"Col"}, Rows: [][]string{{"Val"}}}
	dt.Init(nil)
	html := dt.Render().String()
	css := dt.RenderCSS().String()

	htmlClasses := filterClasses(extractHTMLClasses(html), "datatable")
	cssClasses := filterClasses(extractCSSClasses(css), "datatable")

	for cls := range cssClasses {
		if !htmlClasses[cls] {
			t.Errorf("DataTable CSS class %q does not exist in rendered HTML", cls)
		}
	}
	for cls := range htmlClasses {
		if !cssClasses[cls] {
			t.Errorf("DataTable HTML class %q is unstyled in CSS", cls)
		}
	}
}
