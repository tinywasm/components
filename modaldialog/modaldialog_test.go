package modaldialog

import (
	"testing"

	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/fmt"
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
	if !Contains(html, "class='modal'") {
		t.Error("expected modal class")
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
	m.visible.Set(false)
	htmlHidden := m.Render().String()
	// Show returns a placeholder node when the condition is false
	if htmlHidden == "" {
		t.Error("expected placeholder node string when not visible (Show condition), not empty string")
	}
	if Contains(htmlHidden, "modal") {
		t.Error("should not contain modal when Visible=false")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) String() string            { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []Component { return nil }
