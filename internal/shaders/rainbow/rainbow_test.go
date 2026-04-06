package rainbow

import (
	"testing"

	"github.com/ryanlewis/tint/internal/core"
)

func renderRow(shader core.Shader, text string, frame int) []core.Output {
	runes := []rune(text)
	out := make([]core.Output, len(runes))
	for i, r := range runes {
		out[i] = shader.Shade(core.Cell{X: i, Y: 0, Char: r, Frame: frame, W: len(runes), H: 1})
	}
	return out
}

func TestRainbow(t *testing.T) {
	out := renderRow(New(), "ABCDE", 0)
	// Each position should produce a different hue.
	for i := 1; i < len(out); i++ {
		if out[i].Fg == out[i-1].Fg {
			t.Errorf("rainbow: positions %d and %d have same color", i-1, i)
		}
	}
}

func TestRainbowSkipsSpaces(t *testing.T) {
	out := renderRow(New(), "A B", 0)
	if out[1].Fg.A != 0 {
		t.Errorf("rainbow should not color spaces, got Fg=%+v", out[1].Fg)
	}
	if out[1].Char != ' ' {
		t.Errorf("rainbow should preserve space char, got %q", out[1].Char)
	}
}

func TestRainbowAnimates(t *testing.T) {
	shader := New()
	out0 := shader.Shade(core.Cell{X: 0, Y: 0, Char: 'A', Frame: 0, W: 1, H: 1})
	// Frame 50 wraps hue to 0 again, 25 doesn't.
	out25 := shader.Shade(core.Cell{X: 0, Y: 0, Char: 'A', Frame: 25, W: 1, H: 1})

	if out0.Fg == out25.Fg {
		t.Error("rainbow should produce different colors on different frames")
	}
}

func TestRainbowSingleChar(t *testing.T) {
	out := renderRow(New(), "X", 0)
	if out[0].Fg.A == 0 {
		t.Error("rainbow should color a single character")
	}
	if out[0].Char != 'X' {
		t.Errorf("rainbow should preserve char, got %q", out[0].Char)
	}
}
