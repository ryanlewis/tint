package fire

import (
	"testing"

	"github.com/ryanlewis/tint/internal/core"
)

// luminance returns a rough perceived luminance for comparison.
func luminance(c core.Color) float64 {
	return 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
}

func renderColumn(shader core.Shader, ch rune, h int, frame int) []core.Output {
	out := make([]core.Output, h)
	for y := 0; y < h; y++ {
		out[y] = shader.Shade(core.Cell{X: 0, Y: y, Char: ch, Frame: frame, W: 1, H: h})
	}
	return out
}

func TestFireVerticalGradient(t *testing.T) {
	// Taller grid - bottom rows should be brighter (warmer) than top.
	out := renderColumn(New(), 'X', 10, 0)

	topLum := luminance(out[0].Fg)
	botLum := luminance(out[9].Fg)

	if botLum <= topLum {
		t.Errorf("fire: bottom should be brighter than top (bot=%f, top=%f)", botLum, topLum)
	}
}

func TestFireSkipsSpaces(t *testing.T) {
	shader := New()
	out := shader.Shade(core.Cell{X: 1, Y: 0, Char: ' ', W: 3, H: 1})
	if out.Fg.A != 0 {
		t.Errorf("fire should not color spaces, got Fg=%+v", out.Fg)
	}
}

func TestFireAnimates(t *testing.T) {
	shader := New()
	differs := false
	for y := 0; y < 3; y++ {
		for x := 0; x < 5; x++ {
			out0 := shader.Shade(core.Cell{X: x, Y: y, Char: 'X', Frame: 0, W: 5, H: 3})
			out50 := shader.Shade(core.Cell{X: x, Y: y, Char: 'X', Frame: 50, W: 5, H: 3})
			if out0.Fg != out50.Fg {
				differs = true
			}
		}
	}
	if !differs {
		t.Error("fire should shimmer across frames")
	}
}

func TestFireSingleRow(t *testing.T) {
	// Single row - should not panic, should still produce color.
	out := renderColumn(New(), 'F', 1, 0)
	if out[0].Fg.A == 0 {
		t.Error("fire should color characters even on a single row")
	}
}
