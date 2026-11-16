package ui

import (
	"os"

	"golang.org/x/term"
)

// ANSI color codes, empty when stdout is not a terminal or NO_COLOR is set.
var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

// IsTTY reports whether stdout is a terminal.
var IsTTY = term.IsTerminal(int(os.Stdout.Fd()))

func init() {
	if !IsTTY || os.Getenv("NO_COLOR") != "" {
		ColorReset, ColorRed, ColorGreen, ColorYellow = "", "", "", ""
		ColorBlue, ColorCyan, ColorGray, ColorBold = "", "", "", ""
	}
}
