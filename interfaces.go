package components

// HTMLRenderer renders the component's HTML structure
type HTMLRenderer interface {
	RenderHTML() string
}

// CSSRenderer provides component-specific CSS
type CSSRenderer interface {
	RenderCSS() string
}

// IconProvider provides SVG icons for the global sprite
type IconProvider interface {
	IconSvg() []map[string]string
}

// DisplayNamer provides a human-readable name for navigation
type DisplayNamer interface {
	DisplayName() string
}
