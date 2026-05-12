package contentcard

import (
	"testing"

	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/fmt"
)

func TestCard_Render(t *testing.T) {
	c := &ContentCard{
		Header: &simpleComponent{html: "Header"},
		Body:   &simpleComponent{html: "Body"},
		Footer: &simpleComponent{html: "Footer"},
	}

	html := c.Render().RenderHTML()

	if !HasPrefix(html, "<div") {
		t.Error("expected div tag")
	}

	if !Contains(html, "class='card'") {
		t.Error("expected card class")
	}

	if !Contains(html, "class='card-header'") {
		t.Error("expected card-header")
	}
	if !Contains(html, "Header") {
		t.Error("expected Header content")
	}

	if !Contains(html, "class='card-body'") {
		t.Error("expected card-body")
	}
	if !Contains(html, "Body") {
		t.Error("expected Body content")
	}

	if !Contains(html, "class='card-footer'") {
		t.Error("expected card-footer")
	}
	if !Contains(html, "Footer") {
		t.Error("expected Footer content")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) RenderHTML() string        { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []Component { return nil }
