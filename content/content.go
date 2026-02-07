package content

import (
	"github.com/tinywasm/dom"
)

// Content helps construct a standard Content layout
type Content struct {
	handlerName string
	title       string

	// Layout Sections
	headerControl string
	sidebar       dom.HTMLRenderer
	content       dom.HTMLRenderer

	// Internal tracking
	handler    any
	components []dom.HTMLRenderer // Track all used components for asset collection
}

func New(h any) *Content {
	b := &Content{
		handler: h,
	}

	// 1. Get handler name
	if name, ok := h.(string); ok {
		b.handlerName = name
	} else if named, ok := h.(interface{ HandlerName() string }); ok {
		b.handlerName = named.HandlerName()
	}

	return b
}

func (b *Content) HandlerName() string {
	return b.handlerName
}

func (b *Content) SetTitle(title string) *Content {
	b.title = title
	return b
}

func (b *Content) Title() string {
	return b.title
}

func (b *Content) track(c dom.HTMLRenderer) {
	if c != nil {
		b.components = append(b.components, c)
	}
}

func (b *Content) TrackedComponents() []dom.HTMLRenderer {
	return b.components
}

// WithHeader adds a component to the header controls area.
func (b *Content) WithHeader(c dom.HTMLRenderer) *Content {
	b.headerControl = c.RenderHTML()
	b.track(c)
	return b
}

// WithSidebar sets the sidebar component (formerly searchList).
func (b *Content) WithSidebar(c dom.HTMLRenderer) *Content {
	b.sidebar = c
	b.track(c)
	return b
}

// WithContent sets the main content component (formerly form/buttons).
func (b *Content) WithContent(c dom.HTMLRenderer) *Content {
	b.content = c
	b.track(c)
	return b
}

func (b *Content) RenderHTML() string {
	// 1. Header
	head := &header.ContentHeader{
		Title:    b.title,
		Controls: b.headerControl,
	}
	b.track(head)

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
	b.track(sContainer)

	// 4. Panel
	p := &panel.Panel{
		ID:      b.handlerName,
		Content: head.RenderHTML() + sContainer.RenderHTML(),
	}
	b.track(p)

	return p.RenderHTML()
}
