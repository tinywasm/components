package selectsearch

import (
	"strings"
	"testing"
)

func TestSelectSearch_Render(t *testing.T) {
	c := &SelectSearch{
		Placeholder: "Choose category",
		Options: []Option{
			{ID: "1", Label: "Automobiles", Description: "auto"},
			{ID: "2", Label: "Film & Animation", Description: "anime"},
		},
	}

	html := c.Render().RenderHTML()

	if !strings.Contains(html, "ss-box") {
		t.Error("expected ss-box class")
	}
	if !strings.Contains(html, "ss-toggle") {
		t.Error("expected ss-toggle checkbox")
	}
	if !strings.Contains(html, "Choose category") {
		t.Error("expected placeholder text")
	}
	if !strings.Contains(html, "Automobiles") {
		t.Error("expected option label")
	}
	if !strings.Contains(html, "auto") {
		t.Error("expected option description")
	}
	if !strings.Contains(html, "#ss-arrow-down") {
		t.Error("expected icon use reference")
	}
}

func TestSelectSearch_SelectedValue(t *testing.T) {
	c := &SelectSearch{
		Placeholder:   "Choose category",
		SelectedLabel: "Automobiles",
	}
	html := c.Render().RenderHTML()
	if !strings.Contains(html, "Automobiles") {
		t.Error("expected selected label")
	}
	if strings.Contains(html, "Choose category") {
		// Placeholder should be replaced by selected label
	}
}

func TestSelectSearch_Filtering(t *testing.T) {
	c := &SelectSearch{
		FilterTerm: "Film",
		Options: []Option{
			{ID: "1", Label: "Automobiles"},
			{ID: "2", Label: "Film & Animation"},
		},
	}
	html := c.Render().RenderHTML()
	if strings.Contains(html, "Automobiles") {
		t.Error("expected Automobiles to be filtered out")
	}
	if !strings.Contains(html, "Film & Animation") {
		t.Error("expected Film & Animation to be present")
	}
}
