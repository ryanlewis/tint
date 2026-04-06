package grid

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wantH int
		wantW int
	}{
		{"empty string", "", 0, 0},
		{"single char", "x", 1, 1},
		{"single line", "hello", 1, 5},
		{"multi line", "hi\nthere", 2, 5},
		{"trailing newline stripped", "abc\n", 1, 3},
		{"multiple trailing newlines", "abc\n\n\n", 1, 3}, // TrimRight strips all trailing \n
		{"uneven lines padded", "ab\nlong line\nx", 3, 9},
		{"only newlines", "\n\n", 0, 0}, // all newlines stripped -> empty -> nil
		{"blank lines in middle", "a\n\nb", 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Parse(tt.input)
			if tt.wantH == 0 {
				if g != nil {
					t.Fatalf("expected nil grid, got %d rows", len(g))
				}
				return
			}
			if len(g) != tt.wantH {
				t.Fatalf("height: got %d, want %d", len(g), tt.wantH)
			}
			for i, row := range g {
				if len(row) != tt.wantW {
					t.Errorf("row %d width: got %d, want %d", i, len(row), tt.wantW)
				}
			}
		})
	}
}

func TestParsePadding(t *testing.T) {
	g := Parse("ab\nc")
	if g[1][1] != ' ' {
		t.Errorf("expected space padding, got %q", g[1][1])
	}
}

func TestPad(t *testing.T) {
	g := Parse("AB\nCD")

	padded := Pad(g, 2)
	if len(padded) != 6 { // 2 + 2 + 2
		t.Fatalf("padded height: got %d, want 6", len(padded))
	}
	if len(padded[0]) != 6 { // 2 + 2 + 2
		t.Fatalf("padded width: got %d, want 6", len(padded[0]))
	}

	// Top margin should be spaces.
	for x := 0; x < 6; x++ {
		if padded[0][x] != ' ' {
			t.Errorf("top margin [0][%d] = %q, want space", x, padded[0][x])
		}
	}
	// Original content is centered.
	if padded[2][2] != 'A' || padded[2][3] != 'B' {
		t.Errorf("row 2 = %v, want [' ',' ','A','B',' ',' ']", string(padded[2]))
	}
	if padded[3][2] != 'C' || padded[3][3] != 'D' {
		t.Errorf("row 3 = %v, want [' ',' ','C','D',' ',' ']", string(padded[3]))
	}
}

func TestPadZero(t *testing.T) {
	g := Parse("X")
	padded := Pad(g, 0)
	if len(padded) != 1 || len(padded[0]) != 1 {
		t.Error("pad=0 should return original grid")
	}
}

func TestPadNil(t *testing.T) {
	padded := Pad(nil, 3)
	if padded != nil {
		t.Error("padding nil grid should return nil")
	}
}

func TestParseUnicode(t *testing.T) {
	// Emoji and CJK characters are single runes but multi-byte.
	g := Parse("A🔥B")
	if len(g) != 1 {
		t.Fatalf("expected 1 row, got %d", len(g))
	}
	if len(g[0]) != 3 {
		t.Fatalf("expected width 3 (rune count), got %d", len(g[0]))
	}
	if g[0][0] != 'A' || g[0][1] != '🔥' || g[0][2] != 'B' {
		t.Errorf("grid = %v", g[0])
	}
}

func TestParseWindowsLineEndings(t *testing.T) {
	g := Parse("ab\r\ncd")
	// \r should be preserved as a character (not stripped).
	// This gives width 3: 'a','b','\r' for row 0 and 'c','d',' ' for row 1.
	if len(g) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(g))
	}
	// Row 0 has "ab\r" = 3 runes, row 1 has "cd" = 2 runes padded to 3.
	if len(g[0]) != 3 {
		t.Errorf("row 0 width: got %d, want 3", len(g[0]))
	}
}

func TestParseAllSpaces(t *testing.T) {
	g := Parse("   \n   ")
	if len(g) != 2 || len(g[0]) != 3 {
		t.Errorf("expected 2x3 grid, got %dx%d", len(g), len(g[0]))
	}
	for y, row := range g {
		for x, ch := range row {
			if ch != ' ' {
				t.Errorf("grid[%d][%d] = %q, want space", y, x, ch)
			}
		}
	}
}

func TestParseSingleColumn(t *testing.T) {
	g := Parse("a\nb\nc")
	if len(g) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(g))
	}
	for _, row := range g {
		if len(row) != 1 {
			t.Errorf("expected width 1, got %d", len(row))
		}
	}
}
