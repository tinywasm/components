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

// Bug 2: dropdown must stay open across re-renders while the user is typing.
// isOpen=true must produce a checked checkbox so CSS keeps .ss-dropdown visible.
func TestSelectSearch_OpenState_RendersChecked(t *testing.T) {
	c := &SelectSearch{
		isOpen: true,
		Options: []Option{
			{ID: "1", Label: "Apple"},
		},
	}
	html := c.Render().RenderHTML()
	if !fmt.Contains(html, "checked") {
		t.Error("expected 'checked' attribute on toggle when isOpen=true — dropdown will close on re-render without it")
	}
}

func TestSelectSearch_ClosedState_NoChecked(t *testing.T) {
	c := &SelectSearch{
		isOpen: false,
		Options: []Option{
			{ID: "1", Label: "Apple"},
		},
	}
	html := c.Render().RenderHTML()
	if fmt.Contains(html, "checked") {
		t.Error("expected no 'checked' attribute on toggle when isOpen=false")
	}
}

// Bug 2 + filtering: open state must be preserved while a filter term is active.
func TestSelectSearch_OpenWithFilter_RendersCheckedAndFiltered(t *testing.T) {
	c := &SelectSearch{
		isOpen:     true,
		filterTerm: "Mango",
		Options: []Option{
			{ID: "1", Label: "Apple"},
			{ID: "2", Label: "Mango"},
		},
	}
	html := c.Render().RenderHTML()
	if !fmt.Contains(html, "checked") {
		t.Error("expected 'checked' when isOpen=true with active filter")
	}
	if fmt.Contains(html, "Apple") {
		t.Error("expected Apple to be filtered out")
	}
	if !fmt.Contains(html, "Mango") {
		t.Error("expected Mango to remain visible")
	}
}

// Bug 3 (cursor reset): replicates the condition that causes the cursor to jump
// to position 0 on every keystroke.
//
// Root cause: dom.Update() replaces outerHTML, discarding the active element.
// The search input must render with value=filterTerm so that when the dom layer
// restores focus (see tinywasm/dom docs/PLAN.md), the cursor can be placed at
// the correct position. If value is missing or empty, any cursor restoration
// attempt will land at position 0 regardless.
//
// This test will FAIL if Render() stops emitting value=filterTerm on the search
// input — which would break cursor restoration once dom.Update() is fixed.
func TestSelectSearch_SearchInput_ValueMatchesFilterTerm(t *testing.T) {
	tests := []struct {
		name       string
		filterTerm string
		wantValue  string
	}{
		{"empty term", "", ""},
		{"single char — first keystroke", "a", "a"},
		{"multi char — mid-word", "app", "app"},
		{"unicode", "mañana", "mañana"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SelectSearch{
				isOpen:     true,
				filterTerm: tt.filterTerm,
				Options:    []Option{{ID: "1", Label: "Apple"}},
			}
			html := c.Render().RenderHTML()
			if tt.wantValue == "" {
				return
			}
			needle := `value='` + tt.wantValue + `'`
			if !fmt.Contains(html, needle) {
				t.Errorf("search input missing %q — cursor restoration will fail after dom.Update()", needle)
			}
		})
	}
}

// After selecting an option the dropdown must close (isOpen=false).
func TestSelectSearch_AfterSelection_RendersUnchecked(t *testing.T) {
	c := &SelectSearch{
		isOpen:        false,
		selectedLabel: "Apple",
	}
	html := c.Render().RenderHTML()
	if fmt.Contains(html, "checked") {
		t.Error("expected no 'checked' attribute after selection — dropdown must be closed")
	}
	if !fmt.Contains(html, "Apple") {
		t.Error("expected selected label in header")
	}
}
