// Command tint applies visual shaders to text read from stdin and writes
// truecolor ANSI to stdout. See tint --help for usage.
package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ryanlewis/tint"
	"golang.org/x/term"
)

// Build-time version information, populated via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// CLI describes the tint command-line interface.
//
// Animation behavior: by default, tint animates when stdout is a terminal
// and renders a single frame otherwise (e.g. when piped to a file or another
// command). --once forces a single frame even on a terminal. --animate forces
// animation even when stdout is not a terminal, for cases where you want to
// stream frames to another tool.
type CLI struct {
	Shader  string           `arg:""               help:"Shader to apply."                                                          placeholder:"SHADER"`
	Once    bool             `                     help:"Render a single frame and exit, even on a terminal."                       xor:"mode"`
	Animate bool             `                     help:"Force animation even when stdout is not a terminal."                       xor:"mode"`
	FPS     int              `default:"10"         help:"Animation frame rate."`
	Timeout int              `                     help:"Stop animation after N seconds (0 = no limit)."`
	Pad     int              `                     help:"Add N cells of margin around the text."`
	Version kong.VersionFlag `                     help:"Show version and exit."`
}

func main() {
	cli := CLI{}
	versionStr := fmt.Sprintf("tint %s (%s, built %s)", version, commit, date)

	kong.Parse(&cli,
		kong.Name("tint"),
		kong.Description("Apply visual shaders to text.\n\nShaders: "+strings.Join(tint.Names(), ", ")),
		kong.Vars{"version": versionStr},
		kong.UsageOnError(),
	)

	if cli.FPS <= 0 {
		fmt.Fprintln(os.Stderr, "--fps must be greater than 0")
		os.Exit(1)
	}

	factory, ok := tint.Get(cli.Shader)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown shader: %q\navailable: %s\n", cli.Shader, strings.Join(tint.Names(), ", "))
		os.Exit(1)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading stdin: %v\n", err)
		os.Exit(1)
	}

	grid := tint.ParseGrid(string(input))
	if grid == nil {
		return
	}
	grid = tint.PadGrid(grid, cli.Pad)

	shader := factory(len(grid[0]), len(grid))

	if shouldAnimate(cli) {
		runAnimation(grid, shader, cli.FPS, cli.Timeout)
		return
	}

	out := tint.Render(grid, shader, 0)
	if err := tint.EncodeANSI(os.Stdout, out); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}
}

// shouldAnimate decides whether to animate based on the flags and whether
// stdout is a terminal. --once and --animate are mutually exclusive overrides
// (enforced by kong via the xor tag); without either flag, tint animates when
// stdout is a TTY and renders one frame otherwise.
func shouldAnimate(cli CLI) bool {
	switch {
	case cli.Once:
		return false
	case cli.Animate:
		return true
	default:
		return term.IsTerminal(int(os.Stdout.Fd())) //nolint:gosec // os.Stdout.Fd() is a small non-negative file descriptor
	}
}

// runAnimation loops rendering frames at the given FPS until Ctrl-C, the
// timeout expires, or a write fails. A timeout of 0 means no limit.
func runAnimation(grid [][]rune, shader tint.Shader, fps, timeoutSec int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Fprint(os.Stdout, tint.ANSIHideCursor)
	defer fmt.Fprint(os.Stdout, tint.ANSIShowCursor+tint.ANSIReset)

	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	var deadline <-chan time.Time
	if timeoutSec > 0 {
		deadline = time.After(time.Duration(timeoutSec) * time.Second)
	}

	frame := 0
	for {
		select {
		case <-sig:
			return
		case <-deadline:
			return
		case <-ticker.C:
			out := tint.Render(grid, shader, frame)
			fmt.Fprint(os.Stdout, tint.ANSIClearScreen)
			if err := tint.EncodeANSI(os.Stdout, out); err != nil {
				return
			}
			frame++
		}
	}
}
