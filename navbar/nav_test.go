package navbar

import (
	"testing"

	. "github.com/tinywasm/fmt"
)

func TestNav_Render(t *testing.T) {
	n := &NavBar{
		Items: []NavBarItem{
			{Label: "Home", Route: "home"},
			{Label: "Users", Route: "users", Icon: "user-icon"},
		},
	}

	html := n.Render().RenderHTML()

	if !HasPrefix(html, "<nav") {
		t.Error("expected nav tag")
	}
	if !Contains(html, "class='nav'") {
		t.Error("expected nav class")
	}

	if !Contains(html, "class='nav-list'") {
		t.Error("expected nav-list")
	}

	// Check items
	if !Contains(html, "href='#home'") {
		t.Error("expected home link")
	}
	if !Contains(html, "Home") {
		t.Error("expected Home label")
	}

	if !Contains(html, "href='#users'") {
		t.Error("expected users link")
	}
	if !Contains(html, "Users") {
		t.Error("expected Users label")
	}

	// Check icon
	if !Contains(html, "href='#user-icon'") {
		t.Error("expected user icon symbol")
	}
}
