package dialog

import (
	"testing"

	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/fmt"
)

func TestModal_Render(t *testing.T) {
	m := &DialogWidget{
		Title:   "My Modal",
		Content: &simpleComponent{html: "<p>Content</p>"},
		Visible: true,
	}

	html := m.Render().String()

	// Check main container
	if !Contains(html, "class='modal'") {
		t.Error("expected modal class")
	}
	if Contains(html, "hidden") {
		t.Error("should not be hidden when Visible=true")
	}

	// Check internal structure
	if !Contains(html, "class='modal-backdrop'") {
		t.Error("expected backdrop")
	}
	if !Contains(html, "class='modal-content'") {
		t.Error("expected content container")
	}
	if !Contains(html, "class='modal-header'") {
		t.Error("expected header")
	}
	if !Contains(html, "<h2>My Modal</h2>") { // H2 factory renders <h2>...</h2>
		t.Error("expected title")
	}
	if !Contains(html, "<p>Content</p>") {
		t.Error("expected content")
	}

	// Test hidden
	m.Visible = false
	htmlHidden := m.Render().String()
	if !Contains(htmlHidden, "class='modal hidden'") {
		t.Error("expected hidden class when Visible=false")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) String() string            { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []Component { return nil }
