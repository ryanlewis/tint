// Package fire provides a warm vertical gradient shader with animated shimmer.
package fire

import (
	"math"

	"github.com/ryanlewis/tint/internal/core"
)

// Fire applies a warm gradient (red -> orange -> yellow) that shimmers over time.
type Fire struct{}

// New returns a Fire shader.
func New() Fire {
	return Fire{}
}

// Shade implements core.Shader.
func (f Fire) Shade(c core.Cell) core.Output {
	if c.Char == ' ' {
		return core.Output{Char: ' '}
	}

	// Vertical gradient: bottom is brighter (hotter).
	// Y=0 is top of text, Y=H-1 is bottom. Higher Y = closer to base = hotter.
	h := math.Max(float64(c.H), 1)
	t := float64(c.Y) / h

	// Shimmer using overlapping sine waves.
	shimmer := 0.5 + 0.5*math.Sin(
		float64(c.X)*0.8+
			float64(c.Y)*1.2+
			float64(c.Frame)*0.15,
	)

	// Blend the intensity.
	intensity := t*0.7 + shimmer*0.3

	// Map intensity to a warm hue: 0 (red) -> 30 (orange) -> 50 (yellow).
	hue := intensity * 50
	lightness := 0.3 + intensity*0.25

	return core.Output{
		Char: c.Char,
		Fg:   core.HSLToColor(hue, 1.0, lightness),
	}
}
