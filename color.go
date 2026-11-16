package hearth

import "os"

// ColorConfig holds the inputs that determine whether ANSI colour is
// active — a direct lift of cairn's internal/format.ColorConfig/
// ShouldColor logic, which already had no cairn-specific concepts in it.
type ColorConfig struct {
	IsTTY      bool // stdout is a terminal
	NoColor    bool // --no-color flag was set
	ForceColor bool // --color flag was set
}

// ShouldColor returns true when ANSI colour should be emitted.
//
// Priority (highest first):
//  1. --no-color flag → always off
//  2. FORCE_COLOR env → always on
//  3. --color flag → always on
//  4. NO_COLOR env → off
//  5. IsTTY → on only if terminal
func ShouldColor(c ColorConfig) bool {
	if c.NoColor {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	if c.ForceColor {
		return true
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return c.IsTTY
}

// Color is one of a small set of ANSI colors the render functions in this
// package use. Unlike cairn's ColorStatus/ColorPriority, nothing here is
// keyed to a field name or a known set of values — hearth has no schema
// to key that off of. See Style for the type-based (not field-based)
// coloring policy RenderBlock uses by default.
type Color int

const (
	ColorNone Color = iota
	ColorDim
	ColorCyan
	ColorGreen
	ColorYellow
	ColorRed
)

const (
	ansiReset  = "\x1b[0m"
	ansiDim    = "\x1b[2m"
	ansiCyan   = "\x1b[36m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiRed    = "\x1b[31m"
)

// Colorize wraps s in c's ANSI code when enabled is true; returns s
// unchanged otherwise.
func Colorize(c Color, s string, enabled bool) string {
	if !enabled || c == ColorNone {
		return s
	}
	var code string
	switch c {
	case ColorDim:
		code = ansiDim
	case ColorCyan:
		code = ansiCyan
	case ColorGreen:
		code = ansiGreen
	case ColorYellow:
		code = ansiYellow
	case ColorRed:
		code = ansiRed
	default:
		return s
	}
	return code + s + ansiReset
}
