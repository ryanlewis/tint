package tint

import (
	"maps"
	"slices"
)

var registry = map[string]ShaderFactory{}

// Register adds a shader factory to the global registry under the given name.
// A later Register call with the same name replaces the previous factory.
func Register(name string, factory ShaderFactory) {
	registry[name] = factory
}

// Get returns the shader factory for the given name, or (nil, false) if no
// factory has been registered under that name.
func Get(name string) (ShaderFactory, bool) {
	f, ok := registry[name]
	return f, ok
}

// Names returns all registered shader names in sorted order.
func Names() []string {
	return slices.Sorted(maps.Keys(registry))
}
