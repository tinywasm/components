package selectsearch

import (
	"testing"

	"github.com/tinywasm/fmt"
)

func TestSelectSearch_Render(t *testing.T) {
	c := &SelectSearch{
		Placeholder: "Choose category",
		Options: []SsOption{
			{ID: "1", Label: "Automobiles", Description: "auto"},
			{ID: "2", Label: "Film & Animation", Description: "anime"},
		},
	}
	c.Init(nil)

	html := c.Render().String()

	if !fmt.Contains(html, "selectsearch") {
		t.Error("expected selectsearch class")
	}
	if !fmt.Contains(html, "selectsearch__toggle") {
		t.Error("expected selectsearch__toggle checkbox")
	}
	if !fmt.Contains(html, "Choose category") {
		t.Error("expected placeholder text")
	}
	// Note: options are in the dropdown, which is wrapped in Show(c.isOpen)
	// Since isOpen is false, options won't be in the initial static HTML
}

func TestSelectSearch_SelectedValue(t *testing.T) {
	c := &SelectSearch{
		Placeholder: "Choose category",
	}
	c.Init(nil)
	c.selectedLabel.Set("Automobiles")

	html := c.Render().String()
	if !fmt.Contains(html, "Automobiles") {
		t.Error("expected selected label")
	}
}

func TestSelectSearch_OpenState_RendersChecked(t *testing.T) {
	c := &SelectSearch{
		Options: []SsOption{
			{ID: "1", Label: "Apple"},
		},
	}
	c.Init(nil)
	c.isOpen.Set(true)

	el := c.Render()
	html := el.String()
	if !fmt.Contains(html, "checked") {
		t.Error("expected 'checked' attribute on toggle when isOpen=true")
	}
}

func TestSelectSearch_Filtering(t *testing.T) {
	c := &SelectSearch{
		Options: []SsOption{
			{ID: "1", Label: "Automobiles"},
			{ID: "2", Label: "Film & Animation"},
		},
	}
	c.Init(nil)
	c.isOpen.Set(true)
	c.query.Set("Film")
	// query.Set doesn't trigger OnChange in standard tests unless manually called or using gotest/WASM
	c.rows.Set(c.buildRows("Film"))

	el := c.Render()
	html := el.String()
	if fmt.Contains(html, "Automobiles") {
		t.Error("expected Automobiles to be filtered out")
	}
}
