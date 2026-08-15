package sitenav

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/image"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/widget"
)

// NameSiteNav is the widget name for sitenav.
const NameSiteNav = widget.Name("sitenav")

const (
	PartBrand   = widget.Part("brand")
	PartLogo    = widget.Part("logo")
	PartToggle  = widget.Part("toggle")
	PartMenu    = widget.Part("menu")
	PartNav     = widget.Part("nav")
	PartLink    = widget.Part("link")
	PartActions = widget.Part("actions")
)

var (
	clsNav        = NameSiteNav.Root()
	clsNavBrand   = NameSiteNav.Class(PartBrand)
	clsNavLogo    = NameSiteNav.Class(PartLogo)
	clsNavToggle  = NameSiteNav.Class(PartToggle)
	clsNavMenu    = NameSiteNav.Class(PartMenu)
	clsNavLink    = NameSiteNav.Class(PartLink)
	clsNavActions = NameSiteNav.Class(PartActions)
)

const (
	iconMenu  = svg.Icon("sitenav-menu")
	iconClose = svg.Icon("sitenav-close")
)

// NavItem represents a single link item in the site navigation bar.
type NavItem struct {
	Label  string
	Href   string
	Active bool
}

// SiteNav provides the main site header with brand logo, navigation links,
// actions, and a responsive mobile toggle menu.
type SiteNav struct {
	Element
	WideLogoSrc    string
	CompactLogoSrc string
	LogoAlt        string
	Links          []NavItem
	Actions        []Component
}

func (sn *SiteNav) WidgetName() widget.Name { return NameSiteNav }
func (sn *SiteNav) WidgetKind() widget.Kind { return widget.Region }

func (sn *SiteNav) Render() *Element {
	header := Header().Set(clsNav.AsAttr())

	// Brand section with responsive <picture> logo
	brand := Div().Set(clsNavBrand.AsAttr())
	pic := image.Picture()
	if sn.WideLogoSrc != "" {
		pic.Child(image.Source(sn.WideLogoSrc, "(min-width: 768px)"))
	}
	src := sn.CompactLogoSrc
	if src == "" {
		src = sn.WideLogoSrc
	}
	pic.Child(image.Img(src, sn.LogoAlt).Class(string(clsNavLogo)).AsElement())
	brand.Child(pic)
	header.Child(brand)

	// Toggle button for mobile navigation
	toggleBtn := Button().Set(clsNavToggle.AsAttr()).
		Attr("aria-controls", "sitenav-menu").
		Attr("aria-expanded", "false").
		Attr("aria-label", "Abrir menú de navegación").
		Child(iconMenu.Render("sitenav-icon-open")).
		Child(iconClose.Render("sitenav-icon-close"))
	header.Child(toggleBtn)

	// Menu container
	menu := Div().Set(clsNavMenu.AsAttr()).ID("sitenav-menu")

	nav := Nav()
	for _, link := range sn.Links {
		a := A(link.Href).Set(clsNavLink.AsAttr()).Text(link.Label)
		if link.Active {
			a.Attr("aria-current", "page")
		}
		nav.Child(a)
	}
	menu.Child(nav)

	if len(sn.Actions) > 0 {
		actions := Div().Set(clsNavActions.AsAttr())
		for _, act := range sn.Actions {
			if act != nil {
				actions.Child(act)
			}
		}
		menu.Child(actions)
	}

	header.Child(menu)
	return header
}
