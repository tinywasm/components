package nav

import (
	"github.com/tinywasm/dom"
)

type Nav struct {
	*dom.Element
	Items []NavItem
}

type NavItem struct {
	Label string
	Route string // e.g., "users", "products"
	Icon  string // optional icon ID
}

func (n *Nav) Render() *dom.Element {
	if n.Element == nil {
		n.Element = &dom.Element{}
	}

	nav := dom.Nav().
		Class("nav")

	ul := dom.Ul().Class("nav-list")
	currentHash := dom.GetHash()

	for _, item := range n.Items {
		isActive := currentHash == "#"+item.Route
		li := dom.Li().Class("nav-item")
		if isActive {
			li.Class("active")
		}

		link := NavLink(item.Label, "#"+item.Route, item.Icon)
		li.Add(link)
		ul.Add(li)
	}

	nav.Add(ul)
	return nav
}

// NavLink creates a navigation link with optional icon.
func NavLink(text, hash, icon string) *dom.Element {
	link := dom.A(hash).
		Text(text).
		On("click", func(e dom.Event) {
			e.PreventDefault()
			dom.SetHash(hash)
		})

	if icon != "" {
		// Add icon using SVG sprite
		link.Add(
			dom.Svg().Class("icon").Add(
				dom.Use().Attr("href", "#"+icon),
			),
		)
	}

	return link
}
