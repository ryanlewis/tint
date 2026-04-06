package tint

import (
	"bytes"
	"strings"
	"testing"
)

// shaderFunc adapts a function to the Shader interface.
type shaderFunc func(Cell) Output

func (f shaderFunc) Shade(c Cell) Output { return f(c) }

// ---------------------------------------------------------------------------
// Render
// ---------------------------------------------------------------------------

func TestRender(t *testing.T) {
	grid := ParseGrid("AB")
	factory, _ := Get("solid")
	shader := factory(2, 1)
	out := Render(grid, shader, 0)

	if len(out) != 1 || len(out[0]) != 2 {
		t.Fatalf("expected 1x2 output, got %dx%d", len(out), len(out[0]))
	}
	if out[0][0].Char != 'A' {
		t.Errorf("char: got %c, want A", out[0][0].Char)
	}
	if out[0][0].Fg != Red() {
		t.Errorf("fg: got %+v, want red", out[0][0].Fg)
	}
}

func TestRenderNilGrid(t *testing.T) {
	noop := shaderFunc(func(c Cell) Output { return Output{Char: c.Char} })
	out := Render(nil, noop, 0)
	if out != nil {
		t.Errorf("expected nil output for nil grid, got %d rows", len(out))
	}
}

func TestRenderEmptyGrid(t *testing.T) {
	noop := shaderFunc(func(c Cell) Output { return Output{Char: c.Char} })
	out := Render([][]rune{}, noop, 0)
	if out != nil {
		t.Errorf("expected nil output for empty grid, got %d rows", len(out))
	}
}

func TestRenderPassesCellFields(t *testing.T) {
	grid := ParseGrid("ab\ncd")

	var cells []Cell
	spy := shaderFunc(func(c Cell) Output {
		cells = append(cells, c)
		return Output{Char: c.Char}
	})

	Render(grid, spy, 42)

	if len(cells) != 4 {
		t.Fatalf("expected 4 cells, got %d", len(cells))
	}

	// Top-left.
	tl := cells[0]
	if tl.X != 0 || tl.Y != 0 || tl.Char != 'a' || tl.Frame != 42 || tl.W != 2 || tl.H != 2 {
		t.Errorf("top-left cell: %+v", tl)
	}
	// Bottom-right.
	br := cells[3]
	if br.X != 1 || br.Y != 1 || br.Char != 'd' || br.Frame != 42 || br.W != 2 || br.H != 2 {
		t.Errorf("bottom-right cell: %+v", br)
	}
}

func TestRender1x1Grid(t *testing.T) {
	grid := ParseGrid("X")
	var got Cell
	spy := shaderFunc(func(c Cell) Output {
		got = c
		return Output{Char: c.Char}
	})
	Render(grid, spy, 0)
	if got.W != 1 || got.H != 1 || got.X != 0 || got.Y != 0 {
		t.Errorf("1x1 cell: %+v", got)
	}
}

func TestRenderLargeFrame(t *testing.T) {
	grid := ParseGrid("A")
	var got Cell
	spy := shaderFunc(func(c Cell) Output {
		got = c
		return Output{Char: c.Char}
	})
	Render(grid, spy, 999999)
	if got.Frame != 999999 {
		t.Errorf("expected frame 999999, got %d", got.Frame)
	}
}

// ---------------------------------------------------------------------------
// Registry — verifies built-in shaders are registered at package load
// ---------------------------------------------------------------------------

func TestRegistry(t *testing.T) {
	names := Names()
	expected := []string{"fire", "matrix", "rainbow", "solid"}
	for _, name := range expected {
		if _, ok := Get(name); !ok {
			t.Errorf("shader %q not registered", name)
		}
	}

	// Names should be sorted.
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("Names() not sorted: %v", names)
			break
		}
	}
}

func TestRegistryUnknown(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Error("Get should return false for unknown shader")
	}
}

func TestRegistryCustom(t *testing.T) {
	// Custom shaders added via Register should be discoverable.
	Register("custom-test", func(_, _ int) Shader {
		return shaderFunc(func(c Cell) Output { return Output{Char: c.Char} })
	})
	defer delete(registry, "custom-test")

	if _, ok := Get("custom-test"); !ok {
		t.Error("Register should make shader available via Get")
	}
}

// ---------------------------------------------------------------------------
// Apply — convenience pipeline: parse → render → encode
// ---------------------------------------------------------------------------

func TestApply(t *testing.T) {
	var buf bytes.Buffer
	err := Apply(&buf, "Hi", "rainbow", 0)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("Apply produced no output")
	}
}

func TestApplyUnknownShader(t *testing.T) {
	var buf bytes.Buffer
	err := Apply(&buf, "Hi", "nonexistent", 0)
	if err == nil {
		t.Error("expected error for unknown shader")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention shader name, got: %v", err)
	}
}

func TestApplyEmptyText(t *testing.T) {
	var buf bytes.Buffer
	err := Apply(&buf, "", "rainbow", 0)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Error("Apply on empty text should produce no output")
	}
}

func TestApplyMultiline(t *testing.T) {
	var buf bytes.Buffer
	err := Apply(&buf, "line1\nline2\nline3", "solid", 0)
	if err != nil {
		t.Fatal(err)
	}

	s := buf.String()
	// Should have 3 lines of output (each ending with reset+newline).
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 output lines, got %d", len(lines))
	}
}

// ---------------------------------------------------------------------------
// End-to-end: full pipeline smoke tests
// ---------------------------------------------------------------------------

func TestEndToEndRainbow(t *testing.T) {
	var buf bytes.Buffer
	err := Apply(&buf, "Hello\nWorld", "rainbow", 0)
	if err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	// Should contain truecolor ANSI codes.
	if !strings.Contains(s, "\x1b[38;2;") {
		t.Error("output should contain truecolor codes")
	}
	// Should contain the actual characters.
	if !strings.Contains(s, "H") || !strings.Contains(s, "W") {
		t.Error("output should contain original characters")
	}
	// Should end with reset.
	if !strings.HasSuffix(strings.TrimRight(s, "\n"), ANSIReset) {
		t.Error("output should end with ANSI reset")
	}
}

func TestEndToEndMatrix(t *testing.T) {
	var buf bytes.Buffer
	// Matrix replaces characters, so the output chars may differ from input.
	err := Apply(&buf, "AAAAA\nAAAAA\nAAAAA\nAAAAA\nAAAAA", "matrix", 10)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("matrix should produce output")
	}
}

func TestEndToEndFire(t *testing.T) {
	var buf bytes.Buffer
	err := Apply(&buf, "FIRE\nFIRE\nFIRE", "fire", 0)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("fire should produce output")
	}
}
