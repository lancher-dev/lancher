package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// SelectOption represents an option in the selection menu
type SelectOption struct {
	Value string
	Label string // If empty, Value is used as label
}

// isTerminal checks if the file descriptor is a terminal
func isTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

// setRawMode sets terminal to raw mode for reading arrow keys
func setRawMode(fd uintptr) (*syscall.Termios, error) {
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	newState := oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	return &oldState, nil
}

// restoreTerminal restores terminal to previous state
func restoreTerminal(fd uintptr, oldState *syscall.Termios) error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(oldState)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

// readKey reads a single key or arrow key sequence
func readKey() (string, error) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return "", err
	}

	if n == 1 {
		return string(buf[0]), nil
	}

	// Arrow keys send escape sequences: ESC [ A/B/C/D
	if n == 3 && buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 65:
			return "up", nil
		case 66:
			return "down", nil
		}
	}

	return string(buf[:n]), nil
}

// selectWithArrows displays interactive menu with arrow key navigation
func selectWithArrows(prompt string, items []SelectOption) (string, error) {
	selected := 0

	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit

	// Display initial menu
	fmt.Printf("%s%s%s\n\n", ColorCyan+ColorBold, prompt, ColorReset)

	fd := os.Stdin.Fd()
	oldState, err := setRawMode(fd)
	if err != nil {
		return "", err
	}
	defer restoreTerminal(fd, oldState)

	// Render function
	render := func() {
		// Move cursor to start of options
		fmt.Printf("\033[%dA", len(items))

		for i, item := range items {
			label := item.Label
			if label == "" {
				label = item.Value
			}

			if i == selected {
				// Highlighted option with background
				fmt.Printf("\r  %s%s▸ %s%s\033[K\n", ColorGreen+ColorBold, "\033[7m", label, ColorReset)
			} else {
				// Normal option
				fmt.Printf("\r    %s\033[K\n", label)
			}
		}
	}

	// Initial render
	for range items {
		fmt.Println()
	}
	render()

	// Event loop
	for {
		key, err := readKey()
		if err != nil {
			return "", err
		}

		switch key {
		case "up":
			if selected > 0 {
				selected--
				render()
			}
		case "down":
			if selected < len(items)-1 {
				selected++
				render()
			}
		case "\n", "\r": // Enter key
			fmt.Println() // Move to next line after selection
			return items[selected].Value, nil
		case "\x03": // Ctrl+C
			fmt.Println()
			return "", fmt.Errorf("cancelled by user")
		case "q", "Q":
			fmt.Println()
			return "", fmt.Errorf("cancelled by user")
		}
	}
}

// selectWithNumbers displays numbered menu with numeric input
func selectWithNumbers(prompt string, items []SelectOption) (string, error) {
	fmt.Printf("%s%s%s\n\n", ColorCyan+ColorBold, prompt, ColorReset)

	// Display numbered options
	for i, item := range items {
		label := item.Label
		if label == "" {
			label = item.Value
		}
		fmt.Printf("  %s%d.%s %s\n", ColorGreen, i+1, ColorReset, label)
	}
	fmt.Println()

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sSelect option (1-%d):%s ", ColorYellow, len(items), ColorReset)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	// Parse selection
	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(items) {
		return "", fmt.Errorf("invalid selection: please choose a number between 1 and %d", len(items))
	}

	return items[selection-1].Value, nil
}

// Select displays an interactive selection menu and returns the selected value
// Tries arrow key navigation first, falls back to numbered input if not available
// options can be either []string or []SelectOption
func Select(prompt string, options interface{}) (string, error) {
	var items []SelectOption

	// Convert input to SelectOption slice
	switch v := options.(type) {
	case []string:
		items = make([]SelectOption, len(v))
		for i, s := range v {
			items[i] = SelectOption{Value: s, Label: s}
		}
	case []SelectOption:
		items = v
	default:
		return "", fmt.Errorf("invalid options type")
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	// Check if stdin is a terminal and supports interactive mode
	if isTerminal(os.Stdin.Fd()) {
		// Try arrow key navigation
		result, err := selectWithArrows(prompt, items)
		if err != nil {
			// If arrow mode fails, fall back to numbered selection
			fmt.Printf("\n%sFalling back to numbered selection...%s\n\n", ColorYellow, ColorReset)
			return selectWithNumbers(prompt, items)
		}
		return result, nil
	}

	// Not a terminal, use numbered selection
	return selectWithNumbers(prompt, items)
}

// PromptString prompts the user for a string input
func PromptString(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s%s%s ", ColorCyan, prompt, ColorReset)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return strings.TrimSpace(input), nil
}

// PromptConfirm prompts the user for a yes/no confirmation
func PromptConfirm(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s%s (y/n):%s ", ColorYellow, prompt, ColorReset)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}
