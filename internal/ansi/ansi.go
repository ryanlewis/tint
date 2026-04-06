// Package ansi writes rendered output cells as truecolor ANSI text.
// It optimizes by only emitting escape codes when fg, bg, or style change.
package ansi

import (
	"io"
	"strconv"
	"unicode/utf8"

	"github.com/ryanlewis/tint/internal/core"
)

// ANSI escape sequences.
const (
	ClearScreen = "\x1b[2J\x1b[H"
	HideCursor  = "\x1b[?25l"
	ShowCursor  = "\x1b[?25h"
	Reset       = "\x1b[0m"
)

// sgrCodes maps Style bits to their ANSI SGR code strings.
// Pre-rendered to avoid any formatting in the hot path.
var sgrCodes = [...]struct {
	bit core.Style
	seq string // e.g. "\x1b[1m"
}{
	{core.Bold, "\x1b[1m"},
	{core.Dim, "\x1b[2m"},
	{core.Italic, "\x1b[3m"},
	{core.Underline, "\x1b[4m"},
	{core.Blink, "\x1b[5m"},
	{core.Reverse, "\x1b[7m"},
	{core.Strikethrough, "\x1b[9m"},
}

// Encode writes a 2D Output grid to w as truecolor ANSI text.
func Encode(w io.Writer, out [][]core.Output) error {
	// Pre-allocate a reasonable buffer to avoid repeated growth.
	buf := make([]byte, 0, 4096)
	for _, row := range out {
		buf = encodeRow(buf, row)
	}
	_, err := w.Write(buf)
	return err
}

// rowState tracks the last-emitted attributes within a single row so that
// Encode only writes new escape sequences when something actually changes.
type rowState struct {
	prevFg, prevBg core.Color
	prevStyle      core.Style
	first          bool
}

func encodeRow(buf []byte, row []core.Output) []byte {
	state := rowState{first: true}
	for _, cell := range row {
		buf = state.emitCell(buf, cell)
	}
	buf = append(buf, Reset...)
	buf = append(buf, '\n')
	return buf
}

// emitCell writes a single cell, prefixing it with a reset + new attributes
// only when any of fg/bg/style has changed since the previous cell.
func (s *rowState) emitCell(buf []byte, cell core.Output) []byte {
	if s.first || cell.Fg != s.prevFg || cell.Bg != s.prevBg || cell.Style != s.prevStyle {
		if !s.first {
			buf = append(buf, Reset...)
		}
		buf = appendAttrs(buf, cell)
		s.prevFg = cell.Fg
		s.prevBg = cell.Bg
		s.prevStyle = cell.Style
	}
	buf = utf8.AppendRune(buf, cell.Char)
	s.first = false
	return buf
}

// appendAttrs appends the SGR style bits followed by fg and bg color sequences
// (if set) for a single cell. It writes nothing for transparent attributes.
func appendAttrs(buf []byte, cell core.Output) []byte {
	if cell.Style != 0 {
		for _, sg := range sgrCodes {
			if cell.Style&sg.bit != 0 {
				buf = append(buf, sg.seq...)
			}
		}
	}
	if cell.Fg.A != 0 {
		buf = appendColorSeq(buf, "38", cell.Fg)
	}
	if cell.Bg.A != 0 {
		buf = appendColorSeq(buf, "48", cell.Bg)
	}
	return buf
}

// appendColorSeq appends e.g. "\x1b[38;2;R;G;Bm" to buf without fmt.
func appendColorSeq(buf []byte, kind string, c core.Color) []byte {
	buf = append(buf, "\x1b["...)
	buf = append(buf, kind...)
	buf = append(buf, ";2;"...)
	buf = strconv.AppendUint(buf, uint64(c.R), 10)
	buf = append(buf, ';')
	buf = strconv.AppendUint(buf, uint64(c.G), 10)
	buf = append(buf, ';')
	buf = strconv.AppendUint(buf, uint64(c.B), 10)
	buf = append(buf, 'm')
	return buf
}
