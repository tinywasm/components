package dom

import (
	upstream "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

// Re-export types from upstream
type BaseComponent = upstream.BaseComponent
type Component = upstream.Component
type Event = upstream.Event
type EventHandler = upstream.EventHandler
type Node = upstream.Node
type Element = upstream.Element
type ViewRenderer = upstream.ViewRenderer

// Helper functions that might be missing in upstream
var Mount = upstream.Mount
var Render = upstream.Render
var Append = upstream.Append
var Hydrate = upstream.Hydrate
var Get = upstream.Get
var Update = upstream.Update
var Unmount = upstream.Unmount
var Log = upstream.Log

// ElementBuilder helps construct Node objects fluently
type ElementBuilder struct {
	node Node
}

func (e *ElementBuilder) ID(id string) *ElementBuilder {
	e.node.Attrs = append(e.node.Attrs, fmt.KeyValue{Key: "id", Value: id})
	return e
}

func (e *ElementBuilder) Class(class string) *ElementBuilder {
	// Check if class attribute already exists and append
	found := false
	for i, attr := range e.node.Attrs {
		if attr.Key == "class" {
			e.node.Attrs[i].Value += " " + class
			found = true
			break
		}
	}
	if !found {
		e.node.Attrs = append(e.node.Attrs, fmt.KeyValue{Key: "class", Value: class})
	}
	return e
}

func (e *ElementBuilder) Attr(key, value string) *ElementBuilder {
	e.node.Attrs = append(e.node.Attrs, fmt.KeyValue{Key: key, Value: value})
	return e
}

func (e *ElementBuilder) Text(text string) *ElementBuilder {
	e.node.Children = append(e.node.Children, text)
	return e
}

func (e *ElementBuilder) Append(child any) *ElementBuilder {
	if child != nil {
		if b, ok := child.(*ElementBuilder); ok {
			e.node.Children = append(e.node.Children, b.ToNode())
		} else {
			e.node.Children = append(e.node.Children, child)
		}
	}
	return e
}

func (e *ElementBuilder) OnClick(handler func(Event)) *ElementBuilder {
	e.node.Events = append(e.node.Events, EventHandler{Name: "click", Handler: handler})
	return e
}

func (e *ElementBuilder) OnInput(handler func(Event)) *ElementBuilder {
	e.node.Events = append(e.node.Events, EventHandler{Name: "input", Handler: handler})
	return e
}

func (e *ElementBuilder) OnSubmit(handler func(Event)) *ElementBuilder {
	e.node.Events = append(e.node.Events, EventHandler{Name: "submit", Handler: handler})
	return e
}

func (e *ElementBuilder) ToNode() Node {
	return e.node
}

// ToComponent converts the builder to a Component
func (e *ElementBuilder) ToComponent() Component {
	return &ElementComponent{Builder: e}
}

// ElementComponent wraps an ElementBuilder as a Component
type ElementComponent struct {
	BaseComponent
	Builder *ElementBuilder
}

func (e *ElementComponent) Render() Node {
	return e.Builder.ToNode()
}

// Factory functions
func Div() *ElementBuilder {
	return &ElementBuilder{node: Node{Tag: "div"}}
}

func Button() *ElementBuilder {
	return &ElementBuilder{node: Node{Tag: "button"}}
}

func Tag(tag string) *ElementBuilder {
	return &ElementBuilder{node: Node{Tag: tag}}
}
