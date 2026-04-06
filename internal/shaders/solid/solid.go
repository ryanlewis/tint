// Package solid provides a shader that applies a single foreground color
// to every non-space character.
package solid

import "github.com/ryanlewis/tint/internal/core"

// Solid applies a single foreground color to all non-space characters.
type Solid struct {
	Fg core.Color
}

// New returns a Solid shader using the default red foreground.
func New() Solid {
	return Solid{Fg: core.Red()}
}

// Shade implements core.Shader.
func (s Solid) Shade(c core.Cell) core.Output {
	out := core.Output{Char: c.Char}
	if c.Char != ' ' {
		out.Fg = s.Fg
	}
	return out
}
