package components

import "github.com/tinywasm/dom"

// Registry stores component factories (optional, for future use)
var registry = make(map[string]func() dom.Component)

// Register adds a component factory to the registry.
func Register(name string, factory func() dom.Component) {
	registry[name] = factory
}

// Get retrieves a component by name.
func Get(name string) dom.Component {
	if factory, ok := registry[name]; ok {
		return factory()
	}
	return nil
}

// List returns a list of registered component names.
func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
