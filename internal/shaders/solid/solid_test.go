package solid

import (
	"testing"

	"github.com/ryanlewis/tint/internal/core"
)

func render(shader core.Shader, text string) []core.Output {
	runes := []rune(text)
	out := make([]core.Output, len(runes))
	for i, r := range runes {
		out[i] = shader.Shade(core.Cell{X: i, Y: 0, Char: r, W: len(runes), H: 1})
	}
	return out
}

func TestSolid(t *testing.T) {
	out := render(Solid{Fg: core.Blue()}, "A B")

	if out[0].Fg != core.Blue() {
		t.Errorf("solid: got %+v, want blue", out[0].Fg)
	}
	if out[1].Fg.A != 0 {
		t.Errorf("solid should not color spaces, got %+v", out[1].Fg)
	}
}

func TestSolidAllSpaces(t *testing.T) {
	out := render(Solid{Fg: core.Red()}, "   ")
	for _, cell := range out {
		if cell.Fg.A != 0 {
			t.Errorf("solid should not color spaces, got %+v", cell.Fg)
		}
	}
}

func TestSolidPreservesChar(t *testing.T) {
	out := render(Solid{Fg: core.Red()}, "XYZ")
	if out[0].Char != 'X' || out[1].Char != 'Y' || out[2].Char != 'Z' {
		t.Error("solid should preserve original characters")
	}
}

func TestSolidStaticAcrossFrames(t *testing.T) {
	shader := Solid{Fg: core.Red()}
	out0 := shader.Shade(core.Cell{X: 0, Y: 0, Char: 'A', Frame: 0, W: 1, H: 1})
	out99 := shader.Shade(core.Cell{X: 0, Y: 0, Char: 'A', Frame: 99, W: 1, H: 1})

	if out0 != out99 {
		t.Error("solid should be identical regardless of frame")
	}
}

func TestNew(t *testing.T) {
	s := New()
	if s.Fg != core.Red() {
		t.Errorf("New() default Fg = %+v, want red", s.Fg)
	}
}
