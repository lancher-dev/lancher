// Package ui holds terminal output helpers and interactive prompts.
package ui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ANSI color codes, empty when stdout is not a terminal or NO_COLOR is set.
var (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

// IsTTY reports whether stdout is a terminal.
var IsTTY = term.IsTerminal(int(os.Stdout.Fd()))

var noColor = os.Getenv("NO_COLOR") != ""

func init() {
	if !IsTTY || noColor {
		ColorReset, ColorGreen, ColorYellow = "", "", ""
		ColorCyan, ColorGray, ColorBold = "", "", ""
	}
}

// PrintError writes err to stderr, in color when stderr is a terminal.
func PrintError(err error) {
	mark := "✗"
	if term.IsTerminal(int(os.Stderr.Fd())) && !noColor {
		mark = "\033[31m✗\033[0m"
	}
	fmt.Fprintf(os.Stderr, "%s %v\n", mark, err)
}

// Clean makes the control characters of s visible, so text read from a
// template cannot move the cursor or hide output.
func Clean(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < 0x20 || (r >= 0x7f && r < 0xa0) {
			fmt.Fprintf(&b, "\\x%02x", r)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
