package modal

import (
	"strings"
	"testing"

	"github.com/tinywasm/dom"
)

func TestModal_Render(t *testing.T) {
	m := &Modal{
		Title:   "My Modal",
		Content: &simpleComponent{html: "<p>Content</p>"},
		Visible: true,
	}

	html := m.Render().RenderHTML()

	// Check main container
	if !strings.Contains(html, "class='modal'") {
		t.Error("expected modal class")
	}
	if strings.Contains(html, "hidden") {
		t.Error("should not be hidden when Visible=true")
	}

	// Check internal structure
	if !strings.Contains(html, "class='modal-backdrop'") {
		t.Error("expected backdrop")
	}
	if !strings.Contains(html, "class='modal-content'") {
		t.Error("expected content container")
	}
	if !strings.Contains(html, "class='modal-header'") {
		t.Error("expected header")
	}
	if !strings.Contains(html, "<h2>My Modal</h2>") { // H2 factory renders <h2>...</h2>
		t.Error("expected title")
	}
	if !strings.Contains(html, "<p>Content</p>") {
		t.Error("expected content")
	}

	// Test hidden
	m.Visible = false
	htmlHidden := m.Render().RenderHTML()
	if !strings.Contains(htmlHidden, "class='modal hidden'") {
		t.Error("expected hidden class when Visible=false")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) RenderHTML() string        { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []dom.Component { return nil }
