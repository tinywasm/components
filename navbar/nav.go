package navbar

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
)

var (
	clsNav     Class = "nav"
	clsNavList Class = "nav-list"
	clsNavItem Class = "nav-item"
	clsIcon    Class = "icon"
	clsActive  Class = "active"
)

type NavBar struct {
	Element
	Items []NavBarItem
}

type NavBarItem struct {
	Label string
	Route string // e.g., "users", "products"
	Icon  string // optional icon ID
}

func (n *NavBar) Render() *Element {
	nav := Nav(clsNav.AsAttr())

	ul := Ul(clsNavList.AsAttr())
	currentHash := GetHash()

	for _, item := range n.Items {
		isActive := currentHash == "#"+item.Route
		li := Li().Add(clsNavItem.AsAttr())
		if isActive {
			li.Add(clsActive.AsAttr())
		}

		link := NavLink(item.Label, "#"+item.Route, item.Icon)
		li.Add(link)
		ul.Add(li)
	}

	nav.Add(ul)
	return nav
}

// NavLink creates a navigation link with optional icon.
func NavLink(text, hash, icon string) *Element {
	link := A(hash).
		Text(text).
		On("click", func(e Event) {
			e.PreventDefault()
			SetHash(hash)
		})

	if icon != "" {
		// Add icon using SVG sprite
		link.Add(
			Svg(clsIcon.AsAttr()).Add(
				Use().Attr("href", "#"+icon),
			),
		)
	}

	return link
}
