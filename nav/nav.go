package nav

import (
	"github.com/tinywasm/components"
	dom "github.com/tinywasm/components/internal/dom"
)

type Nav struct {
	components.BaseComponent
	Items []NavItem
}

type NavItem struct {
	Label string
	Route string // e.g., "users", "products"
	Icon  string // optional icon ID
}

func (n *Nav) Render() dom.Node {
	nav := dom.Tag("nav").
		ID(n.ID()).
		Class("nav")

	ul := dom.Tag("ul").Class("nav-list")

	for _, item := range n.Items {
		li := dom.Tag("li").Class("nav-item")

		link := dom.Tag("a").
			Attr("href", "#"+item.Route).
			Text(item.Label)

		if item.Icon != "" {
			// Add icon using SVG sprite
			link.Append(
				dom.Tag("svg").Class("icon").Append(
					dom.Tag("use").Attr("href", "#"+item.Icon),
				),
			)
		}

		li.Append(link)
		ul.Append(li)
	}

	nav.Append(ul)
	return nav.ToNode()
}
