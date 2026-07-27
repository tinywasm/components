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

	html := c.Render().String()

	if !HasPrefix(html, "<div") {
		t.Error("expected div tag")
	}

	if !Contains(html, "class='contentcard'") {
		t.Error("expected contentcard class")
	}

	if !Contains(html, "class='contentcard__header'") {
		t.Error("expected contentcard__header")
	}
	if !Contains(html, "Header") {
		t.Error("expected Header content")
	}

	if !Contains(html, "class='contentcard__body'") {
		t.Error("expected contentcard__body")
	}
	if !Contains(html, "Body") {
		t.Error("expected Body content")
	}

	if !Contains(html, "class='contentcard__footer'") {
		t.Error("expected contentcard__footer")
	}
	if !Contains(html, "Footer") {
		t.Error("expected Footer content")
	}
}

type simpleComponent struct {
	html string
}

func (s *simpleComponent) String() string            { return s.html }
func (s *simpleComponent) GetID() string             { return "" }
func (s *simpleComponent) SetID(id string)           {}
func (s *simpleComponent) Children() []Component { return nil }
