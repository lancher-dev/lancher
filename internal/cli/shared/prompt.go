package shared

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// stdin buffers standard input for all prompts.
var stdin = bufio.NewReader(os.Stdin)

func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// rawMode puts the terminal in raw mode and returns the function restoring it.
func rawMode() (func(), error) {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return func() { term.Restore(fd, state) }, nil
}

// readLine reads a line from stdin and returns def when it is empty.
func readLine(def string) (string, error) {
	line, err := stdin.ReadString('\n')
	if err != nil && (err != io.EOF || line == "") {
		return "", err
	}
	if line = strings.TrimSpace(line); line == "" {
		return def, nil
	}
	return line, nil
}

// readKey reads one key in raw mode: a printable character, or one of
// "enter", "backspace", "ctrl+c", "up", "down", "left", "right" and "esc".
func readKey() (string, error) {
	r, _, err := stdin.ReadRune()
	if err != nil {
		return "", err
	}
	switch r {
	case '\r', '\n':
		return "enter", nil
	case 127, 8:
		return "backspace", nil
	case 3:
		return "ctrl+c", nil
	case 27:
		return readEscape()
	}
	if r < 32 {
		return "", nil
	}
	return string(r), nil
}

// readEscape consumes an escape sequence already started by ESC.
func readEscape() (string, error) {
	if stdin.Buffered() == 0 {
		return "esc", nil
	}
	if b, _ := stdin.Peek(1); b[0] != '[' {
		return "esc", nil
	}
	stdin.ReadByte()
	var final byte
	for stdin.Buffered() > 0 {
		b, err := stdin.ReadByte()
		if err != nil {
			return "", err
		}
		if b >= '@' && b <= '~' {
			final = b
			break
		}
	}
	switch final {
	case 'A':
		return "up", nil
	case 'B':
		return "down", nil
	case 'C':
		return "right", nil
	case 'D':
		return "left", nil
	}
	return "", nil
}

// PromptStringWithDefault asks for a line of text, showing def as a gray
// placeholder that is used when the answer is empty.
func PromptStringWithDefault(prompt, def string) (string, error) {
	fmt.Printf("%s%s%s ", ColorCyan, prompt, ColorReset)

	var restore func()
	if stdinIsTerminal() {
		restore, _ = rawMode()
	}
	if restore == nil {
		fmt.Printf("(default: %s%s%s) ", ColorGray, def, ColorReset)
		return readLine(def)
	}
	defer restore()

	placeholder := func() {
		if def != "" {
			fmt.Printf("%s%s%s\033[%dD", ColorGray, def, ColorReset, utf8.RuneCountInString(def))
		}
	}
	placeholder()

	var input []rune
	for {
		key, err := readKey()
		if err != nil {
			return "", err
		}
		switch key {
		case "enter":
			fmt.Print("\r\033[K")
			if len(input) == 0 {
				return def, nil
			}
			return string(input), nil
		case "ctrl+c":
			fmt.Print("\r\033[K")
			return "", ErrCancelled
		case "backspace":
			if len(input) == 0 {
				continue
			}
			input = input[:len(input)-1]
			fmt.Print("\033[D\033[K")
			if len(input) == 0 {
				placeholder()
			}
		default:
			if utf8.RuneCountInString(key) != 1 {
				continue
			}
			if len(input) == 0 {
				fmt.Print("\033[K")
			}
			input = append(input, []rune(key)...)
			fmt.Print(key)
		}
	}
}

// PromptConfirmWithDefault asks a yes/no question; Enter picks def.
func PromptConfirmWithDefault(prompt string, def bool) (bool, error) {
	hint, defChar := "y/N", "n"
	if def {
		hint, defChar = "Y/n", "y"
	}
	fmt.Printf("%s%s (%s):%s ", ColorYellow, prompt, hint, ColorReset)

	var restore func()
	if stdinIsTerminal() {
		restore, _ = rawMode()
	}
	if restore == nil {
		answer, err := readLine(defChar)
		if err != nil {
			return false, err
		}
		answer = strings.ToLower(answer)
		return answer == "y" || answer == "yes", nil
	}
	defer restore()

	fmt.Printf("%s%s%s\033[D", ColorGray, defChar, ColorReset)
	for {
		key, err := readKey()
		if err != nil {
			return false, err
		}
		switch strings.ToLower(key) {
		case "enter":
			fmt.Printf("\033[K%s\r\n", defChar)
			return def, nil
		case "y":
			fmt.Print("\033[Ky\r\n")
			return true, nil
		case "n":
			fmt.Print("\033[Kn\r\n")
			return false, nil
		case "ctrl+c":
			fmt.Print("\033[K\r\n")
			return false, ErrCancelled
		}
	}
}
