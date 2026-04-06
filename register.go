package tint

import (
	"github.com/ryanlewis/tint/internal/shaders/fire"
	"github.com/ryanlewis/tint/internal/shaders/matrix"
	"github.com/ryanlewis/tint/internal/shaders/rainbow"
	"github.com/ryanlewis/tint/internal/shaders/solid"
)

// Built-in shaders are registered at package load. Custom shaders can be
// added at runtime via Register.
//
//nolint:gochecknoinits // registering built-in shaders is the purpose of this file
func init() {
	Register("solid", func(_, _ int) Shader { return solid.New() })
	Register("fire", func(_, _ int) Shader { return fire.New() })
	Register("rainbow", func(_, _ int) Shader { return rainbow.New() })
	Register("matrix", func(w, h int) Shader { return matrix.New(w, h) })
}
