package core

import "math"

// RGB creates an opaque Color.
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

// RGBA creates a Color with explicit alpha.
func RGBA(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// HSLToColor converts HSL values to an opaque Color.
// h is in degrees [0, 360), s and l are in [0, 1].
func HSLToColor(h, s, l float64) Color {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return RGB(
		uint8(math.Round((r+m)*255)),
		uint8(math.Round((g+m)*255)),
		uint8(math.Round((b+m)*255)),
	)
}

// White returns the named color white as an opaque Color.
func White() Color { return RGB(255, 255, 255) }

// Black returns the named color black as an opaque Color.
func Black() Color { return RGB(0, 0, 0) }

// Red returns the named color red as an opaque Color.
func Red() Color { return RGB(255, 0, 0) }

// Green returns the named color green as an opaque Color.
func Green() Color { return RGB(0, 255, 0) }

// Blue returns the named color blue as an opaque Color.
func Blue() Color { return RGB(0, 0, 255) }

// Yellow returns the named color yellow as an opaque Color.
func Yellow() Color { return RGB(255, 255, 0) }

// Cyan returns the named color cyan as an opaque Color.
func Cyan() Color { return RGB(0, 255, 255) }

// Magenta returns the named color magenta as an opaque Color.
func Magenta() Color { return RGB(255, 0, 255) }

// Orange returns the named color orange as an opaque Color.
func Orange() Color { return RGB(255, 165, 0) }
