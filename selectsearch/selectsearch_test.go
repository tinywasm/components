package selectsearch

import (
	"github.com/tinywasm/fmt"
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

	if !fmt.Contains(html, "ss-box") {
		t.Error("expected ss-box class")
	}
	if !fmt.Contains(html, "ss-toggle") {
		t.Error("expected ss-toggle checkbox")
	}
	if !fmt.Contains(html, "Choose category") {
		t.Error("expected placeholder text")
	}
	if !fmt.Contains(html, "Automobiles") {
		t.Error("expected option label")
	}
	if !fmt.Contains(html, "auto") {
		t.Error("expected option description")
	}
	if !fmt.Contains(html, "#ss-arrow-down") {
		t.Error("expected icon use reference")
	}
}

func TestSelectSearch_SelectedValue(t *testing.T) {
	c := &SelectSearch{
		Placeholder:   "Choose category",
		selectedLabel: "Automobiles",
	}
	html := c.Render().RenderHTML()
	if !fmt.Contains(html, "Automobiles") {
		t.Error("expected selected label")
	}
	if fmt.Contains(html, "Choose category") {
		// Placeholder should be replaced by selected label
	}
}

func TestSelectSearch_Filtering(t *testing.T) {
	c := &SelectSearch{
		filterTerm: "Film",
		Options: []Option{
			{ID: "1", Label: "Automobiles"},
			{ID: "2", Label: "Film & Animation"},
		},
	}
	html := c.Render().RenderHTML()
	if fmt.Contains(html, "Automobiles") {
		t.Error("expected Automobiles to be filtered out")
	}
	if !fmt.Contains(html, "Film & Animation") {
		t.Error("expected Film & Animation to be present")
	}
}
