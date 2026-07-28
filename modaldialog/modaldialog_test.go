package modaldialog

import (
	"regexp"
	"strings"
	"testing"

	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/fmt"
)

var (
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
)

func TestModal_Render(t *testing.T) {
	m := &ModalDialog{
		Title:   "My Modal",
		Content: &simpleComponent{html: "<p>Content</p>"},
	}
	m.Init(nil)
	m.visible.Set(true)

	html := m.Render().String()

	// Check main container
	if !Contains(html, "class='modaldialog'") {
		t.Error("expected modaldialog class")
	}

	// Check internal structure
	if !Contains(html, "class='modaldialog__backdrop'") {
		t.Error("expected backdrop")
	}
	if !Contains(html, "class='modaldialog__panel'") {
		t.Error("expected content container")
	}
	if !Contains(html, "class='modaldialog__header'") {
		t.Error("expected header")
	}
	if !Contains(html, "<h2>My Modal</h2>") { // H2 factory renders <h2>...</h2>
		t.Error("expected title")
	}
	if !Contains(html, "<p>Content</p>") {
		t.Error("expected content")
	}

	// Test hidden
	m.visible.Set(false)
	htmlHidden := m.Render().String()
	// Show returns a placeholder node when the condition is false
	if htmlHidden == "" {
		t.Error("expected placeholder node string when not visible (Show condition), not empty string")
	}
	if Contains(htmlHidden, "modaldialog") {
		t.Error("should not contain modaldialog when Visible=false")
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

	md := &ModalDialog{Title: "Title", Content: &simpleComponent{html: "<p>Content</p>"}}
	md.Init(nil)
	md.Open()
	html := md.Render().String()
	css := md.Style().Stylesheet().String()

	htmlClasses := filterClasses(extractHTMLClasses(html), "modaldialog")
	cssClasses := filterClasses(extractCSSClasses(css), "modaldialog")

	for cls := range cssClasses {
		if !htmlClasses[cls] {
			t.Errorf("ModalDialog CSS class %q does not exist in rendered HTML", cls)
		}
	}
	for cls := range htmlClasses {
		if !cssClasses[cls] {
			t.Errorf("ModalDialog HTML class %q is unstyled in CSS", cls)
		}
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) String() string            { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []Component { return nil }
