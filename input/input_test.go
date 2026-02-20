package input

import (
	"strings"
	"testing"
)

func TestInput_Render(t *testing.T) {
	inp := &Input{
		Label:       "Username",
		Placeholder: "Enter username",
		Value:       "test",
		Type:        "text",
		Required:    true,
	}

	html := inp.Render().RenderHTML()

	if !strings.Contains(html, "class='input-group'") {
		t.Error("expected input-group class")
	}

	if !strings.Contains(html, "<label") {
		t.Error("expected label tag")
	}
	if !strings.Contains(html, "Username") {
		t.Error("expected label text 'Username'")
	}

	if !strings.Contains(html, "<input") {
		t.Error("expected input tag")
	}
	if !strings.Contains(html, "type='text'") {
		t.Error("expected type='text'")
	}
	if !strings.Contains(html, "placeholder='Enter username'") {
		t.Error("expected placeholder")
	}
	if !strings.Contains(html, "value='test'") {
		t.Error("expected value='test'")
	}
	if !strings.Contains(html, "required=''") {
		t.Error("expected required attribute")
	}
}
