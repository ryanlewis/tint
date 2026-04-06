// Package tint applies visual shaders — color, animation, character
// substitution — to 2D text grids and emits truecolor ANSI.
//
// A shader decides how each cell of a text grid should be rendered: its
// character, its foreground and background colors, and any style bits. See
// Shader for the interface a custom shader must implement.
package tint

import (
	"fmt"
	"io"

	"github.com/ryanlewis/tint/internal/ansi"
	"github.com/ryanlewis/tint/internal/core"
	"github.com/ryanlewis/tint/internal/grid"
)

// Core types. These are aliases for the equivalent types in internal/core,
// which allows the internal packages to share the types without an import
// cycle while keeping the public API stable.
type (
	// Cell is the input to a shader — one pixel of the text grid.
	Cell = core.Cell
	// Color is an RGBA color. A=0 means transparent (no color applied).
	Color = core.Color
	// Output is the visual properties a shader decides for a single pixel.
	Output = core.Output
	// Style is a bitfield of text attributes, composable with |.
	Style = core.Style
	// Shader applies a visual effect to each cell of a text grid.
	Shader = core.Shader
	// ShaderFactory constructs a shader for a grid of the given dimensions.
	ShaderFactory = core.ShaderFactory
)

// Style bits map to ANSI SGR codes.
const (
	Bold          = core.Bold          // ANSI SGR 1
	Dim           = core.Dim           // ANSI SGR 2
	Italic        = core.Italic        // ANSI SGR 3
	Underline     = core.Underline     // ANSI SGR 4
	Blink         = core.Blink         // ANSI SGR 5
	Reverse       = core.Reverse       // ANSI SGR 7
	Strikethrough = core.Strikethrough // ANSI SGR 9
)

// ANSI escape sequences useful when driving an animation loop directly.
const (
	ANSIClearScreen = ansi.ClearScreen
	ANSIHideCursor  = ansi.HideCursor
	ANSIShowCursor  = ansi.ShowCursor
	ANSIReset       = ansi.Reset
)

// RGB creates an opaque Color.
func RGB(r, g, b uint8) Color { return core.RGB(r, g, b) }

// RGBA creates a Color with explicit alpha.
func RGBA(r, g, b, a uint8) Color { return core.RGBA(r, g, b, a) }

// HSLToColor converts HSL values to an opaque Color.
// h is in degrees [0, 360), s and l are in [0, 1].
func HSLToColor(h, s, l float64) Color { return core.HSLToColor(h, s, l) }

// White returns the named color white as an opaque Color.
func White() Color { return core.White() }

// Black returns the named color black as an opaque Color.
func Black() Color { return core.Black() }

// Red returns the named color red as an opaque Color.
func Red() Color { return core.Red() }

// Green returns the named color green as an opaque Color.
func Green() Color { return core.Green() }

// Blue returns the named color blue as an opaque Color.
func Blue() Color { return core.Blue() }

// Yellow returns the named color yellow as an opaque Color.
func Yellow() Color { return core.Yellow() }

// Cyan returns the named color cyan as an opaque Color.
func Cyan() Color { return core.Cyan() }

// Magenta returns the named color magenta as an opaque Color.
func Magenta() Color { return core.Magenta() }

// Orange returns the named color orange as an opaque Color.
func Orange() Color { return core.Orange() }

// ParseGrid converts a multiline string into a 2D rune grid.
// Lines are split on newlines and padded with spaces to uniform width.
// A trailing newline is stripped if present.
func ParseGrid(text string) [][]rune { return grid.Parse(text) }

// PadGrid adds a margin of space characters around a grid.
// The margin is applied uniformly on all four sides.
func PadGrid(g [][]rune, margin int) [][]rune { return grid.Pad(g, margin) }

// EncodeANSI writes a 2D Output grid to w as truecolor ANSI text.
// It emits escape codes only when fg, bg, or style actually change.
func EncodeANSI(w io.Writer, out [][]Output) error { return ansi.Encode(w, out) }

// Render applies a shader to every cell of the grid for the given frame,
// returning a 2D grid of Output values. It returns nil for an empty grid.
func Render(g [][]rune, s Shader, frame int) [][]Output {
	if len(g) == 0 {
		return nil
	}
	h := len(g)
	w := len(g[0])

	out := make([][]Output, h)
	for y := range g {
		row := make([]Output, w)
		for x, ch := range g[y] {
			row[x] = s.Shade(Cell{
				X:     x,
				Y:     y,
				Char:  ch,
				Frame: frame,
				W:     w,
				H:     h,
			})
		}
		out[y] = row
	}
	return out
}

// Apply is a convenience function that parses text into a grid, looks up the
// named shader, renders a single frame, and encodes the result as ANSI.
func Apply(w io.Writer, text, shaderName string, frame int) error {
	factory, ok := Get(shaderName)
	if !ok {
		return fmt.Errorf("unknown shader: %q", shaderName)
	}

	g := ParseGrid(text)
	if g == nil {
		return nil
	}

	shader := factory(len(g[0]), len(g))
	out := Render(g, shader, frame)
	return EncodeANSI(w, out)
}
