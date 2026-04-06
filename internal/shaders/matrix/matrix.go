// Package matrix provides a parallax "matrix rain" shader — katakana rain
// falling on multiple depth planes, interacting with the text in the center.
package matrix

import (
	"math/rand/v2"

	"github.com/ryanlewis/tint/internal/core"
)

// Precomputed constant color for the matrix shader hot path.
var matrixDim = core.HSLToColor(120, 0.8, 0.12)

// dimBright is the base brightness for text at rest.
const dimBright = 0.12

// Matrix is a stateful parallax matrix rain shader.
type Matrix struct {
	layers    []rainLayer
	height    int
	lastFrame int
	rng       *rand.Rand
}

// rainLayer is a single depth plane of falling rain drops.
type rainLayer struct {
	drops    []int // per-column Y position of the drop head (-big = inactive)
	speeds   []int // per-column: advance every N frames
	ticks    []int // per-column frame counter
	trailLen int
	bright   float64 // 0–1: peak brightness at head
	gap      int     // min gap before a drop restarts (controls density)
	isFront  bool    // renders in front of text (can overlay characters)
}

// New constructs a Matrix shader sized for a w×h grid.
func New(w, h int) *Matrix {
	//nolint:gosec // animation randomness, not security-sensitive
	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	trail := h
	if trail < 6 {
		trail = 6
	}

	m := &Matrix{
		height:    h,
		lastFrame: -1,
		rng:       rng,
		layers: []rainLayer{
			// Background plane — very sparse distant rain behind the text.
			newLayer(rng, w, h, 4, 3, 4, 0.10, h*8, false), // dim, short trails, long gaps

			// Text sweep plane — drives the color on actual text characters.
			newLayer(rng, w, h, 2, 2, trail, 0.55, h, false), // mid

			// Foreground planes — rare bright katakana streaking across text.
			newLayer(rng, w, h, 1, 2, 3, 0.45, h*8, true),  // near fg
			newLayer(rng, w, h, 1, 1, 2, 0.65, h*12, true), // far fg: very rare
		},
	}
	return m
}

func newLayer(rng *rand.Rand, w, h, baseSpeed, speedVar, trailLen int, bright float64, gap int, front bool) rainLayer {
	drops := make([]int, w)
	speeds := make([]int, w)
	ticks := make([]int, w)
	for i := range drops {
		// Stagger starts widely — sparse layers start most columns inactive.
		drops[i] = -(rng.IntN(gap) + h)
		speeds[i] = baseSpeed + rng.IntN(max(speedVar, 1))
	}
	return rainLayer{
		drops:    drops,
		speeds:   speeds,
		ticks:    ticks,
		trailLen: trailLen,
		bright:   bright,
		gap:      gap,
		isFront:  front,
	}
}

// Shade implements core.Shader.
func (m *Matrix) Shade(c core.Cell) core.Output {
	if c.Frame != m.lastFrame {
		m.advance()
		m.lastFrame = c.Frame
	}

	isSpace := c.Char == ' '

	// Foreground layers may fully determine the output for spaces.
	fgGlow, fgHead, early, matched := m.scanForeground(c, isSpace)
	if matched {
		return early
	}

	if !isSpace {
		return m.shadeText(c, fgGlow, fgHead)
	}
	return m.shadeBackgroundSpace(c)
}

// scanForeground iterates the foreground layers. For spaces at a head or in a
// katakana trail, it returns an early output (matched=true). For text cells,
// it returns the glow intensity and whether any foreground head touches the
// cell (matched=false).
func (m *Matrix) scanForeground(c core.Cell, isSpace bool) (glow float64, head bool, early core.Output, matched bool) {
	for i := len(m.layers) - 1; i >= 0; i-- {
		ly := &m.layers[i]
		if !ly.isFront {
			continue
		}
		dropHead := ly.drops[c.X]
		dist := dropHead - c.Y
		glowLen := ly.trailLen * 4

		switch {
		case dist == 0:
			if isSpace {
				return 0, false, core.Output{
					Char:  randomKatakana(m.rng),
					Fg:    core.HSLToColor(120, 0.9, ly.bright),
					Style: core.Bold,
				}, true
			}
			head = true
			glow = ly.bright
		case dist > 0 && dist < glowLen:
			if isSpace && dist < ly.trailLen {
				fade := 1.0 - float64(dist)/float64(ly.trailLen)
				return 0, false, core.Output{
					Char: randomKatakana(m.rng),
					Fg:   core.HSLToColor(120, 1.0, fade*ly.bright*0.5),
				}, true
			}
			fade := 1.0 - float64(dist)/float64(glowLen)
			if g := fade * ly.bright; g > glow {
				glow = g
			}
		}
	}
	return glow, head, core.Output{}, false
}

// shadeText colors an actual text character using the sweep layer brightness,
// with optional boost from foreground layer glow.
func (m *Matrix) shadeText(c core.Cell, fgGlow float64, fgHead bool) core.Output {
	sweep := &m.layers[1]
	sweepHead := sweep.drops[c.X]
	sweepDist := sweepHead - c.Y

	baseBright := sweepBrightness(sweepDist, sweep.trailLen)

	// Foreground glow boosts brightness on text characters.
	if fgGlow > 0 {
		glowBright := 0.15 + fgGlow*0.5
		if glowBright > baseBright {
			baseBright = glowBright
		}
	}

	if fgHead || sweepDist == 0 {
		return core.Output{Char: c.Char, Fg: core.HSLToColor(120, 0.6, baseBright), Style: core.Bold}
	}
	if baseBright <= dimBright {
		return core.Output{Char: c.Char, Fg: matrixDim, Style: core.Dim}
	}
	return core.Output{Char: c.Char, Fg: core.HSLToColor(120, 1.0, baseBright)}
}

// sweepBrightness returns the base brightness for a text cell given its
// distance from the text-sweep drop head.
func sweepBrightness(sweepDist, trailLen int) float64 {
	switch {
	case sweepDist == 0:
		return 1.0
	case sweepDist > 0 && sweepDist < trailLen:
		return dimBright + (1.0-float64(sweepDist)/float64(trailLen))*0.42
	default:
		return dimBright
	}
}

// shadeBackgroundSpace renders a space cell using the background rain layer.
func (m *Matrix) shadeBackgroundSpace(c core.Cell) core.Output {
	bg := &m.layers[0]
	bgHead := bg.drops[c.X]
	bgDist := bgHead - c.Y
	if bgDist >= 0 && bgDist < bg.trailLen {
		fade := 1.0 - float64(bgDist)/float64(bg.trailLen)
		return core.Output{
			Char: randomKatakana(m.rng),
			Fg:   core.HSLToColor(120, 1.0, fade*bg.bright*0.4),
		}
	}
	return core.Output{Char: ' '}
}

func (m *Matrix) advance() {
	for li := range m.layers {
		ly := &m.layers[li]
		for i := range ly.drops {
			ly.ticks[i]++
			if ly.ticks[i] >= ly.speeds[i] {
				ly.ticks[i] = 0
				ly.drops[i]++
				if ly.drops[i] > m.height+ly.trailLen {
					// Large gap before restarting — controls density.
					ly.drops[i] = -(m.rng.IntN(ly.gap) + 1)
					ly.speeds[i] = max(ly.speeds[i]-1+m.rng.IntN(3), 1)
				}
			}
		}
	}
}

// randomKatakana returns a random half-width katakana character (U+FF66–U+FF9F).
// Half-width variants are single-cell width in terminals, unlike full-width
// katakana (U+30A0–U+30FF) which occupy two columns and break grid alignment.
func randomKatakana(rng *rand.Rand) rune {
	//nolint:gosec // result is bounded to [0xFF66, 0xFF9F], well within rune range
	return rune(0xFF66 + rng.IntN(0xFF9F-0xFF66+1))
}
