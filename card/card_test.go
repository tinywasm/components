package card

import (
	"strings"
	"testing"

	"github.com/tinywasm/dom"
)

func TestCard_Render(t *testing.T) {
	c := &Card{
		Header: &simpleComponent{html: "Header"},
		Body:   &simpleComponent{html: "Body"},
		Footer: &simpleComponent{html: "Footer"},
	}

	html := c.Render().RenderHTML()

	if !strings.HasPrefix(html, "<div") {
		t.Error("expected div tag")
	}

	if !strings.Contains(html, "class='card'") {
		t.Error("expected card class")
	}

	if !strings.Contains(html, "class='card-header'") {
		t.Error("expected card-header")
	}
	if !strings.Contains(html, "Header") {
		t.Error("expected Header content")
	}

	if !strings.Contains(html, "class='card-body'") {
		t.Error("expected card-body")
	}
	if !strings.Contains(html, "Body") {
		t.Error("expected Body content")
	}

	if !strings.Contains(html, "class='card-footer'") {
		t.Error("expected card-footer")
	}
	if !strings.Contains(html, "Footer") {
		t.Error("expected Footer content")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) RenderHTML() string        { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []dom.Component { return nil }
