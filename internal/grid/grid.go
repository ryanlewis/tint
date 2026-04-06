// Package grid converts text into a rectangular rune grid and provides
// uniform padding around such grids.
package grid

import "strings"

// Pad adds a margin of space characters around a grid.
// The margin is applied uniformly on all four sides.
func Pad(grid [][]rune, margin int) [][]rune {
	if grid == nil || margin <= 0 {
		return grid
	}

	oldH := len(grid)
	oldW := len(grid[0])
	newW := oldW + margin*2
	newH := oldH + margin*2

	padded := make([][]rune, newH)
	for y := range padded {
		row := make([]rune, newW)
		for x := range row {
			row[x] = ' '
		}
		// Copy original content into the center.
		if y >= margin && y < margin+oldH {
			copy(row[margin:], grid[y-margin])
		}
		padded[y] = row
	}
	return padded
}

// Parse converts a multiline string into a 2D rune grid.
// Lines are split on newlines and padded with spaces to uniform width.
// A trailing newline is stripped if present.
func Parse(text string) [][]rune {
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return nil
	}

	lines := strings.Split(text, "\n")

	// Find the maximum line width.
	maxWidth := 0
	for _, line := range lines {
		runes := []rune(line)
		if len(runes) > maxWidth {
			maxWidth = len(runes)
		}
	}

	grid := make([][]rune, len(lines))
	for i, line := range lines {
		runes := []rune(line)
		row := make([]rune, maxWidth)
		copy(row, runes)
		for j := len(runes); j < maxWidth; j++ {
			row[j] = ' '
		}
		grid[i] = row
	}

	return grid
}
