package content

import (
	"github.com/tinywasm/dom"
)

// Content helps construct a standard Content layout
type Content struct {
	title string

	// Layout Sections
	header  dom.Component
	sidebar dom.Component
	content dom.Component
}

func New(title string) *Content {
	b := &Content{
		title: title,
	}
	return b
}

// WithHeader adds a component to the header controls area.
func (b *Content) WithHeader(c dom.Component) *Content {
	b.header = c
	return b
}

// WithSidebar sets the sidebar component (formerly searchList).
func (b *Content) WithSidebar(c dom.Component) *Content {
	b.sidebar = c
	return b
}

// WithContent sets the main content component (formerly form/buttons).
func (b *Content) WithContent(c dom.Component) *Content {
	b.content = c
	return b
}

func (b *Content) RenderHTML() string {
	// 1. Header

	var headerHtml string
	if b.header != nil {
		headerHtml = b.header.RenderHTML()
	}

	// 2. Sidebar + Content Layout
	// Combine sidebar and content into the expected layout structure
	var contentHtml string

	if b.sidebar != nil {
		contentHtml += b.sidebar.RenderHTML()
	}

	if b.content != nil {
		contentHtml += b.content.RenderHTML()
	}

	// 3. Scroll Container (wraps sidebar + content)
	sContainer := &container.ScrollContainer{
		Content: contentHtml,
	}

	// 4. Panel
	p := &panel.Panel{
		ID:      b.handlerName,
		Content: head.RenderHTML() + sContainer.RenderHTML(),
	}

	return p.RenderHTML()
}
