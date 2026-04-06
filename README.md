# tint

Apply visual shaders — colour, animation, character substitution — to text in
your terminal. Tint is a small Go library and CLI that treats a block of text
as a 2D grid of cells and lets a "shader" decide how each cell should be
rendered as truecolor ANSI.

## CLI

```sh
echo "HELLO" | tint rainbow
figlet "FIRE" | tint fire --fps 30
echo "HELLO" | tint matrix --once > frame.ans
```

By default, tint animates when stdout is a terminal (looping until Ctrl-C or
`--timeout`) and renders a single frame otherwise. That means piping to a
file, pager, or another command gives you one clean frame without extra
flags, while running interactively Just Works.

Built-in shaders:

| Shader    | Description                                              |
| --------- | -------------------------------------------------------- |
| `rainbow` | Horizontal hue gradient that shifts across frames        |
| `fire`    | Warm vertical gradient with animated shimmer             |
| `matrix`  | Parallax katakana rain on multiple depth planes          |
| `solid`   | A single foreground colour                               |

Flags:

```
  --once          render a single frame and exit, even on a terminal
  --animate       force animation even when stdout is not a terminal
                  (useful for streaming frames to another tool)
  --fps N         animation frame rate (default: 10)
  --timeout N     stop animation after N seconds (default: no limit)
  --pad N         add N cells of margin around the text
```

`--once` and `--animate` are mutually exclusive.

## Library

```go
import "github.com/ryanlewis/tint"

func main() {
    _ = tint.Apply(os.Stdout, "HELLO\nWORLD", "rainbow", 0)
}
```

A shader is anything that implements the `Shader` interface:

```go
type Shader interface {
    Shade(Cell) Output
}
```

`Cell` carries the rune, its grid coordinates, the grid dimensions and the
current animation frame. `Output` is what the shader decides for that cell:
the (possibly substituted) rune, foreground and background colours, and a
style bitfield (bold, italic, underline, etc.).

Register a custom shader once and it becomes available to both the library
and the CLI:

```go
tint.Register("invert", func(w, h int) tint.Shader {
    return invertShader{}
})
```

## Install

```sh
go install github.com/ryanlewis/tint/cmd/tint@latest
```

Or build from source:

```sh
just build
```

## Development

```sh
just            # list targets
just test       # run tests
just lint       # run golangci-lint
just ci         # lint + test + build
```

## License

MIT
