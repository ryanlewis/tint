// Package rainbow provides a horizontal rainbow shader.
package rainbow

import (
	"math"

	"github.com/ryanlewis/tint/internal/core"
)

// Rainbow applies a horizontal hue gradient that shifts with animation frames.
type Rainbow struct{}

// New returns a Rainbow shader.
func New() Rainbow {
	return Rainbow{}
}

// Shade implements core.Shader.
func (r Rainbow) Shade(c core.Cell) core.Output {
	if c.Char == ' ' {
		return core.Output{Char: ' '}
	}

	t := float64(c.X)/math.Max(float64(c.W), 1) + float64(c.Frame)*0.02
	hue := math.Mod(t*360, 360)

	return core.Output{
		Char: c.Char,
		Fg:   core.HSLToColor(hue, 1.0, 0.5),
	}
}
