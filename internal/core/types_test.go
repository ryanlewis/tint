package core

import "testing"

func TestStyleBitfield(t *testing.T) {
	s := Bold | Italic | Underline
	if s&Bold == 0 {
		t.Error("Bold not set")
	}
	if s&Italic == 0 {
		t.Error("Italic not set")
	}
	if s&Underline == 0 {
		t.Error("Underline not set")
	}
	if s&Dim != 0 {
		t.Error("Dim should not be set")
	}
	if s&Blink != 0 {
		t.Error("Blink should not be set")
	}
}

func TestStyleAllBits(t *testing.T) {
	all := Bold | Dim | Italic | Underline | Blink | Reverse | Strikethrough
	for _, bit := range []Style{Bold, Dim, Italic, Underline, Blink, Reverse, Strikethrough} {
		if all&bit == 0 {
			t.Errorf("style bit %d should be set", bit)
		}
	}
}

func TestStyleZeroValue(t *testing.T) {
	var s Style
	if s != 0 {
		t.Error("zero Style should have no bits set")
	}
}
