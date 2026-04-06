package core

import "testing"

func TestHSLToColor(t *testing.T) {
	tests := []struct {
		name                string
		h, s, l             float64
		wantR, wantG, wantB uint8
	}{
		{"pure red", 0, 1.0, 0.5, 255, 0, 0},
		{"pure green", 120, 1.0, 0.5, 0, 255, 0},
		{"pure blue", 240, 1.0, 0.5, 0, 0, 255},
		{"yellow", 60, 1.0, 0.5, 255, 255, 0},
		{"cyan", 180, 1.0, 0.5, 0, 255, 255},
		{"magenta", 300, 1.0, 0.5, 255, 0, 255},
		{"white", 0, 0, 1.0, 255, 255, 255},
		{"black", 0, 0, 0, 0, 0, 0},
		{"mid grey", 0, 0, 0.5, 128, 128, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := HSLToColor(tt.h, tt.s, tt.l)
			if c.R != tt.wantR || c.G != tt.wantG || c.B != tt.wantB {
				t.Errorf("HSL(%v,%v,%v) = RGB(%d,%d,%d), want RGB(%d,%d,%d)",
					tt.h, tt.s, tt.l, c.R, c.G, c.B, tt.wantR, tt.wantG, tt.wantB)
			}
			if c.A != 255 {
				t.Errorf("alpha should be 255, got %d", c.A)
			}
		})
	}
}

func TestHSLHueWrapping(t *testing.T) {
	// Negative hue wraps.
	if HSLToColor(-120, 1.0, 0.5) != HSLToColor(240, 1.0, 0.5) {
		t.Error("-120 should equal 240")
	}
	// Hue > 360 wraps.
	if HSLToColor(480, 1.0, 0.5) != HSLToColor(120, 1.0, 0.5) {
		t.Error("480 should equal 120")
	}
	// Exactly 360 wraps to 0.
	if HSLToColor(360, 1.0, 0.5) != HSLToColor(0, 1.0, 0.5) {
		t.Error("360 should equal 0")
	}
}

func TestHSLBoundaries(t *testing.T) {
	// All six 60-degree sector boundaries should not panic.
	for hue := 0.0; hue < 360; hue += 59.9 {
		c := HSLToColor(hue, 1.0, 0.5)
		if c.A != 255 {
			t.Errorf("hue=%v: alpha should be 255", hue)
		}
	}
}

func TestRGBConstructor(t *testing.T) {
	c := RGB(10, 20, 30)
	if c.R != 10 || c.G != 20 || c.B != 30 || c.A != 255 {
		t.Errorf("RGB(10,20,30) = %+v", c)
	}
}

func TestRGBAConstructor(t *testing.T) {
	c := RGBA(10, 20, 30, 0)
	if c.A != 0 {
		t.Errorf("RGBA with A=0: got A=%d", c.A)
	}
}

func TestNamedColors(t *testing.T) {
	tests := []struct {
		name    string
		got     Color
		r, g, b uint8
	}{
		{"White", White(), 255, 255, 255},
		{"Black", Black(), 0, 0, 0},
		{"Red", Red(), 255, 0, 0},
		{"Green", Green(), 0, 255, 0},
		{"Blue", Blue(), 0, 0, 255},
		{"Yellow", Yellow(), 255, 255, 0},
		{"Cyan", Cyan(), 0, 255, 255},
		{"Magenta", Magenta(), 255, 0, 255},
		{"Orange", Orange(), 255, 165, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.R != tt.r || tt.got.G != tt.g || tt.got.B != tt.b || tt.got.A != 255 {
				t.Errorf("%s() = %+v, want RGB(%d,%d,%d,255)", tt.name, tt.got, tt.r, tt.g, tt.b)
			}
		})
	}
}
