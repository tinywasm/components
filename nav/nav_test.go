package nav

import (
	"strings"
	"testing"
)

func TestNav_Render(t *testing.T) {
	n := &Nav{
		Items: []NavItem{
			{Label: "Home", Route: "home"},
			{Label: "Users", Route: "users", Icon: "user-icon"},
		},
	}

	html := n.Render().RenderHTML()

	if !strings.HasPrefix(html, "<nav") {
		t.Error("expected nav tag")
	}
	if !strings.Contains(html, "class='nav'") {
		t.Error("expected nav class")
	}

	if !strings.Contains(html, "class='nav-list'") {
		t.Error("expected nav-list")
	}

	// Check items
	if !strings.Contains(html, "href='#home'") {
		t.Error("expected home link")
	}
	if !strings.Contains(html, "Home") {
		t.Error("expected Home label")
	}

	if !strings.Contains(html, "href='#users'") {
		t.Error("expected users link")
	}
	if !strings.Contains(html, "Users") {
		t.Error("expected Users label")
	}

	// Check icon
	if !strings.Contains(html, "href='#user-icon'") {
		t.Error("expected user icon symbol")
	}
}
