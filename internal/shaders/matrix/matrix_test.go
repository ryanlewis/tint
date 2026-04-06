package matrix

import (
	"testing"

	"github.com/ryanlewis/tint/internal/core"
)

func render(shader core.Shader, w, h, frame int, fill rune) [][]core.Output {
	out := make([][]core.Output, h)
	for y := 0; y < h; y++ {
		row := make([]core.Output, w)
		for x := 0; x < w; x++ {
			row[x] = shader.Shade(core.Cell{X: x, Y: y, Char: fill, Frame: frame, W: w, H: h})
		}
		out[y] = row
	}
	return out
}

func renderText(shader core.Shader, lines []string, frame int) [][]core.Output {
	h := len(lines)
	w := len([]rune(lines[0]))
	out := make([][]core.Output, h)
	for y := 0; y < h; y++ {
		runes := []rune(lines[y])
		row := make([]core.Output, w)
		for x := 0; x < w; x++ {
			row[x] = shader.Shade(core.Cell{X: x, Y: y, Char: runes[x], Frame: frame, W: w, H: h})
		}
		out[y] = row
	}
	return out
}

func TestMatrixStateful(t *testing.T) {
	shader := New(5, 5)
	// Run through many frames and check that colors change over time.
	differs := false
	prev := render(shader, 5, 5, 0, '.')
	for f := 1; f <= 50; f++ {
		cur := render(shader, 5, 5, f, '.')
		for y := range prev {
			for x := range prev[y] {
				if prev[y][x].Fg != cur[y][x].Fg {
					differs = true
				}
			}
		}
		prev = cur
	}
	if !differs {
		t.Error("matrix should produce different output across frames")
	}
}

func TestMatrixKatakanaInSpaces(t *testing.T) {
	// Katakana rain appears in spaces (background + foreground layers).
	// Use a grid of spaces so rain has somewhere to land.
	shader := New(12, 12)

	hasKatakana := false
	for f := 0; f < 300; f++ {
		out := render(shader, 12, 12, f, ' ')
		for _, row := range out {
			for _, cell := range row {
				if cell.Char >= 0xFF66 && cell.Char <= 0xFF9F {
					hasKatakana = true
				}
			}
		}
	}
	if !hasKatakana {
		t.Error("matrix should show katakana rain in spaces")
	}
}

func TestMatrixPreservesText(t *testing.T) {
	// Most characters should be preserved (only occasional head flicker).
	shader := New(5, 5)
	lines := []string{"HELLO", "WORLD", "HELLO", "WORLD", "HELLO"}

	out := renderText(shader, lines, 0)

	preserved := 0
	total := 0
	for y, row := range out {
		runes := []rune(lines[y])
		for x, cell := range row {
			if runes[x] != ' ' {
				total++
				if cell.Char == runes[x] {
					preserved++
				}
			}
		}
	}
	// At least 70% of non-space chars should be preserved.
	ratio := float64(preserved) / float64(total)
	if ratio < 0.7 {
		t.Errorf("matrix should mostly preserve text, but only %.0f%% preserved", ratio*100)
	}
}

func TestMatrixHeadIsBold(t *testing.T) {
	shader := New(10, 10)

	// Run frames until we see bold cells from different layers.
	foundWhiteHead := false // text sweep head
	foundGreenHead := false // foreground layer head
	for f := 0; f < 100; f++ {
		out := render(shader, 10, 10, f, 'X')
		for _, row := range out {
			for _, cell := range row {
				if cell.Style&core.Bold != 0 {
					if cell.Fg == core.White() {
						foundWhiteHead = true
					} else if cell.Fg.G > cell.Fg.R { // green-ish
						foundGreenHead = true
					}
				}
			}
		}
	}
	if !foundWhiteHead {
		t.Error("matrix should produce white bold heads (text sweep layer)")
	}
	if !foundGreenHead {
		t.Error("matrix should produce green bold heads (foreground layers)")
	}
}

func TestMatrixWidth1(t *testing.T) {
	shader := New(1, 5)
	// Should not panic with a single column.
	for f := 0; f < 20; f++ {
		render(shader, 1, 5, f, '.')
	}
}

func TestMatrixFrameIdempotent(t *testing.T) {
	// Calling Shade twice with the same frame should not advance drops.
	// The Style (bold = head, dim = background) pattern should be identical.
	m := New(5, 5)

	out1 := render(m, 5, 5, 0, '.')
	out2 := render(m, 5, 5, 0, '.') // same frame again

	for y := range out1 {
		for x := range out1[y] {
			// Bold indicates head position — this should not change.
			bold1 := out1[y][x].Style&core.Bold != 0
			bold2 := out2[y][x].Style&core.Bold != 0
			if bold1 != bold2 {
				t.Fatalf("cell [%d][%d]: head position changed between same-frame renders", y, x)
			}
			// Dim indicates background — should also match.
			dim1 := out1[y][x].Style&core.Dim != 0
			dim2 := out2[y][x].Style&core.Dim != 0
			if dim1 != dim2 {
				t.Fatalf("cell [%d][%d]: dim state changed between same-frame renders", y, x)
			}
		}
	}
}

func TestMatrixDropsCycle(t *testing.T) {
	// Run many frames - drops should eventually cycle (no out-of-bounds panic).
	shader := New(3, 3)
	for f := 0; f < 500; f++ {
		render(shader, 3, 3, f, '.')
	}
}
