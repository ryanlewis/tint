package ansi

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ryanlewis/tint/internal/core"
)

func TestEncode(t *testing.T) {
	out := [][]core.Output{
		{
			{Char: 'H', Fg: core.Red()},
			{Char: 'i', Fg: core.Red()},
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}

	s := buf.String()
	if !strings.Contains(s, "\x1b[38;2;255;0;0m") {
		t.Error("missing red fg escape code")
	}
	if !strings.Contains(s, "H") || !strings.Contains(s, "i") {
		t.Error("missing character output")
	}
	if !strings.HasSuffix(s, Reset+"\n") {
		t.Errorf("should end with reset+newline, got suffix: %q", s[len(s)-10:])
	}
}

func TestEncodeTransparent(t *testing.T) {
	out := [][]core.Output{
		{{Char: ' ', Fg: core.Color{}, Bg: core.Color{}}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	if strings.Contains(s, "38;2") || strings.Contains(s, "48;2") {
		t.Errorf("transparent color should not emit codes, got: %q", s)
	}
}

func TestEncodeDiffOptimization(t *testing.T) {
	out := [][]core.Output{
		{
			{Char: 'A', Fg: core.Red()},
			{Char: 'B', Fg: core.Red()},
			{Char: 'C', Fg: core.Red()},
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	count := strings.Count(s, "\x1b[38;2;255;0;0m")
	if count != 1 {
		t.Errorf("expected 1 red fg code (diff opt), got %d in: %q", count, s)
	}
}

func TestEncodeColorChange(t *testing.T) {
	out := [][]core.Output{
		{
			{Char: 'R', Fg: core.Red()},
			{Char: 'G', Fg: core.Green()},
			{Char: 'B', Fg: core.Blue()},
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	if !strings.Contains(s, "\x1b[38;2;255;0;0m") {
		t.Error("missing red")
	}
	if !strings.Contains(s, "\x1b[38;2;0;255;0m") {
		t.Error("missing green")
	}
	if !strings.Contains(s, "\x1b[38;2;0;0;255m") {
		t.Error("missing blue")
	}
}

func TestEncodeBackground(t *testing.T) {
	out := [][]core.Output{
		{{Char: 'X', Bg: core.Blue()}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	if !strings.Contains(s, "\x1b[48;2;0;0;255m") {
		t.Errorf("missing blue bg code in: %q", s)
	}
}

func TestEncodeFgTransparentBgSet(t *testing.T) {
	// Fg transparent, Bg set - should emit bg but not fg.
	out := [][]core.Output{
		{{Char: 'X', Fg: core.Color{}, Bg: core.Red()}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	if strings.Contains(s, "38;2") {
		t.Error("should not emit fg code when transparent")
	}
	if !strings.Contains(s, "48;2;255;0;0") {
		t.Error("should emit bg code")
	}
}

func TestEncodeStyles(t *testing.T) {
	out := [][]core.Output{
		{{Char: 'B', Style: core.Bold}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	if !strings.Contains(s, "\x1b[1m") {
		t.Errorf("missing bold SGR code in: %q", s)
	}
}

func TestEncodeCombinedStyles(t *testing.T) {
	out := [][]core.Output{
		{{Char: 'X', Style: core.Bold | core.Italic | core.Underline, Fg: core.Red()}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	if !strings.Contains(s, "\x1b[1m") {
		t.Error("missing bold")
	}
	if !strings.Contains(s, "\x1b[3m") {
		t.Error("missing italic")
	}
	if !strings.Contains(s, "\x1b[4m") {
		t.Error("missing underline")
	}
}

func TestEncodeResetPerRow(t *testing.T) {
	out := [][]core.Output{
		{{Char: 'A', Fg: core.Red()}},
		{{Char: 'B', Fg: core.Green()}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	// Each row should end with a reset.
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if !strings.HasSuffix(line, Reset) {
			t.Errorf("line %d should end with reset: %q", i, line)
		}
	}
}

func TestEncodeEmptyGrid(t *testing.T) {
	var buf bytes.Buffer
	err := Encode(&buf, nil)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("nil grid should produce no output, got %d bytes", buf.Len())
	}
}

func TestEncodeCharacterSubstitution(t *testing.T) {
	// Shader that replaces char - e.g. matrix outputs katakana.
	out := [][]core.Output{
		{{Char: 'ア', Fg: core.Green()}},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	if !strings.Contains(s, "ア") {
		t.Error("should contain substituted character")
	}
}

func TestEncodeDiffResetBetweenDifferentCells(t *testing.T) {
	// When color changes, a reset should be emitted before the new codes.
	out := [][]core.Output{
		{
			{Char: 'A', Fg: core.Red()},
			{Char: 'B', Fg: core.Blue()},
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, out); err != nil {
		t.Fatal(err)
	}
	s := buf.String()

	// After 'A', there should be a reset before the blue code.
	redIdx := strings.Index(s, "\x1b[38;2;255;0;0m")
	blueIdx := strings.Index(s, "\x1b[38;2;0;0;255m")
	resetBetween := strings.Index(s[redIdx+1:], Reset)

	if resetBetween < 0 || redIdx+1+resetBetween > blueIdx {
		t.Error("expected reset between color changes")
	}
}
