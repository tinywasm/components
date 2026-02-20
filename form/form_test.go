package form

import (
	"strings"
	"testing"

	"github.com/tinywasm/dom"
)

func TestForm_Render(t *testing.T) {
	f := &Form{
		Fields: []dom.Component{
			&simpleComponent{html: "<input>"},
		},
	}

	html := f.Render().RenderHTML()

	if !strings.HasPrefix(html, "<form") {
		t.Error("expected form tag")
	}

	if !strings.Contains(html, "class='form'") {
		t.Error("expected form class")
	}

	if !strings.Contains(html, "<input>") {
		t.Error("expected input field")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) RenderHTML() string        { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []dom.Component { return nil }
