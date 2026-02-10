package components

import "github.com/tinywasm/dom"

// BaseComponent provides common utilities for component authors
// (Currently just an alias, but can be extended later)
type BaseComponent = dom.BaseComponent

// Props pattern (optional utility)
type Props map[string]any

// Get returns the value for a key as an interface{}.
func (p Props) Get(key string) any {
	return p[key]
}

// String returns the value for a key as a string, or an empty string if not found or not a string.
func (p Props) String(key string) string {
	if v, ok := p[key].(string); ok {
		return v
	}
	return ""
}

// Int returns the value for a key as an int, or 0 if not found or not an int.
func (p Props) Int(key string) int {
	if v, ok := p[key].(int); ok {
		return v
	}
	return 0
}
