// Package core defines the fundamental types shared between the public tint
// package and its internal implementation packages. The public tint package
// re-exports these via type aliases so callers see them as tint.Cell, etc.
package core

// Cell is the input to a shader — one pixel of the text grid.
type Cell struct {
	X, Y  int  // position in the grid
	Char  rune // the character at this position
	Frame int  // animation frame (0 for static)
	W, H  int  // grid dimensions (for normalising coordinates)
}

// Color is an RGBA color. A=0 means transparent (no color applied).
type Color struct {
	R, G, B, A uint8
}

// Style is a bitfield of text attributes, composable with |.
type Style uint8

// Style bits map to ANSI SGR codes.
const (
	Bold          Style = 1 << iota // ANSI SGR 1
	Dim                             // ANSI SGR 2
	Italic                          // ANSI SGR 3
	Underline                       // ANSI SGR 4
	Blink                           // ANSI SGR 5
	Reverse                         // ANSI SGR 7
	Strikethrough                   // ANSI SGR 9
)

// Output is the visual properties a shader decides for a single pixel.
type Output struct {
	Char  rune  // allows character substitution
	Fg    Color // foreground color (A=0 = no override)
	Bg    Color // background color (A=0 = no background)
	Style Style // zero value = no style applied
}

// Shader applies a visual effect to each cell of a text grid.
type Shader interface {
	Shade(c Cell) Output
}

// ShaderFactory constructs a shader for a grid of the given dimensions.
type ShaderFactory func(w, h int) Shader
