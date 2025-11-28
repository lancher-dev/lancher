package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// SelectOption represents an option in the selection menu
type SelectOption struct {
	Value string
	Label string // If empty, Value is used as label
}

// isTerminal checks if the file descriptor is a terminal
func isTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// Terminal state for restoration
var savedState *term.State

// setRawMode sets terminal to raw mode for reading arrow keys
func setRawMode(fd int) error {
	state, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	savedState = state
	return nil
}

// restoreTerminal restores terminal to previous state
func restoreTerminal(fd int) error {
	if savedState == nil {
		return nil
	}
	return term.Restore(fd, savedState)
}

// readKey reads a single key or arrow key sequence
func readKey() (string, error) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return "", err
	}

	// Check for escape sequences (arrow keys)
	if n == 3 && buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 65:
			return "up", nil
		case 66:
			return "down", nil
		case 67:
			return "right", nil
		case 68:
			return "left", nil
		}
	}

	// Check for single characters
	if n == 1 {
		switch buf[0] {
		case 10, 13: // Enter
			return "enter", nil
		case 3: // Ctrl+C
			return "ctrl+c", nil
		case 113: // q
			return "q", nil
		}
	}

	return string(buf[:n]), nil
}

// selectWithArrows provides interactive selection with arrow key navigation
func selectWithArrows(prompt string, options []SelectOption) (string, error) {
	fd := int(os.Stdin.Fd())
	selected := 0

	// Set terminal to raw mode
	if err := setRawMode(fd); err != nil {
		// Fall back to numbered selection if raw mode fails
		return selectWithNumbers(prompt, options)
	}
	defer restoreTerminal(fd)

	// Hide cursor and show initial selection
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	renderOptions := func() {
		fmt.Printf("\r\033[K%s%s%s\n", ColorCyan, prompt, ColorReset)
		for i, opt := range options {
			if i == selected {
				fmt.Printf("\r\033[K%s>%s %s\n", ColorGreen, ColorReset, opt.Label)
			} else {
				fmt.Printf("\r\033[K%s•%s %s\n", ColorGray, ColorReset, opt.Label)
			}
		}
		fmt.Printf("\033[%dA", len(options)+1)
	}

	renderOptions()

	for {
		key, err := readKey()
		if err != nil {
			return "", err
		}

		switch key {
		case "up":
			if selected > 0 {
				selected--
				renderOptions()
			}
		case "down":
			if selected < len(options)-1 {
				selected++
				renderOptions()
			}
		case "enter":
			// Clear the display
			for i := 0; i <= len(options); i++ {
				fmt.Print("\r\033[K\n")
			}
			fmt.Printf("\033[%dA", len(options)+1)
			return options[selected].Value, nil
		case "ctrl+c", "q":
			fmt.Print("\033[?25h")
			fmt.Println("\r\033[K\nCancelled")
			return "", fmt.Errorf("cancelled by user")
		}
	}
}

// selectWithNumbers provides numbered selection as fallback
func selectWithNumbers(prompt string, options []SelectOption) (string, error) {
	fmt.Println(prompt)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt.Label)
	}
	fmt.Print("Enter number: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("invalid selection")
	}

	return options[choice-1].Value, nil
}

// Select prompts user to select from options with arrow key navigation
// Falls back to numbered selection if terminal doesn't support raw mode
func Select(prompt string, choices []string) (string, error) {
	if len(choices) == 0 {
		return "", fmt.Errorf("no choices provided")
	}

	// Convert string choices to SelectOption
	options := make([]SelectOption, len(choices))
	for i, choice := range choices {
		options[i] = SelectOption{
			Value: choice,
			Label: choice,
		}
	}

	fd := int(os.Stdin.Fd())
	if isTerminal(fd) {
		return selectWithArrows(prompt, options)
	}

	return selectWithNumbers(prompt, options)
}

// SelectWithOptions prompts user to select from SelectOption with custom labels
func SelectWithOptions(prompt string, options []SelectOption) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	// Set default labels if empty
	for i := range options {
		if options[i].Label == "" {
			options[i].Label = options[i].Value
		}
	}

	fd := int(os.Stdin.Fd())
	if isTerminal(fd) {
		return selectWithArrows(prompt, options)
	}

	return selectWithNumbers(prompt, options)
}

// PromptString prompts for a text input
func PromptString(prompt string) (string, error) {
	fmt.Printf("%s%s%s ", ColorCyan, prompt, ColorReset)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// Clear the prompt line (go up one line and clear it)
	fmt.Print("\033[1A\033[K")
	return strings.TrimSpace(input), nil
}

// PromptStringWithDefault prompts for text input with a default value shown as gray placeholder
func PromptStringWithDefault(prompt, defaultValue string) (string, error) {
	fmt.Printf("%s%s%s ", ColorCyan, prompt, ColorReset)

	// Check if terminal supports raw mode for interactive placeholder
	fd := int(os.Stdin.Fd())
	if !isTerminal(fd) {
		// Fallback: simple prompt without placeholder
		fmt.Printf("(default: %s%s%s) ", ColorGray, defaultValue, ColorReset)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if input == "" {
			return defaultValue, nil
		}
		return input, nil
	}

	// Set terminal to raw mode for character-by-character reading
	state, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback on error
		fmt.Printf("(default: %s%s%s) ", ColorGray, defaultValue, ColorReset)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if input == "" {
			return defaultValue, nil
		}
		return input, nil
	}
	defer term.Restore(fd, state)

	// Print placeholder in gray at cursor position
	fmt.Printf("%s%s%s", ColorGray, defaultValue, ColorReset)
	// Move cursor back to start of placeholder
	for i := 0; i < len(defaultValue); i++ {
		fmt.Print("\033[D")
	}

	var inputBuffer []rune
	placeholderVisible := true

	buf := make([]byte, 1)
	for {
		// Read single byte
		n, err := os.Stdin.Read(buf)
		if err != nil {
			term.Restore(fd, state)
			return "", err
		}
		if n == 0 {
			continue
		}

		b := buf[0]

		// Handle Enter key
		if b == 13 || b == 10 {
			fmt.Println() // Move to next line
			break
		}

		// Handle Backspace/Delete
		if b == 127 || b == 8 {
			if len(inputBuffer) > 0 {
				// Remove last character from buffer
				inputBuffer = inputBuffer[:len(inputBuffer)-1]

				if len(inputBuffer) == 0 {
					// Input is now empty - show placeholder again
					// Move back one position, then clear to end of line and show placeholder
					fmt.Print("\033[D")  // Move cursor back
					fmt.Print("\033[K")  // Clear to end of line
					fmt.Printf("%s%s%s", ColorGray, defaultValue, ColorReset)
					// Move cursor back to start of placeholder
					for i := 0; i < len(defaultValue); i++ {
						fmt.Print("\033[D")
					}
					placeholderVisible = true
				} else {
					// Move back and clear character
					fmt.Print("\033[D \033[D")
				}
			}
			continue
		}

		// Handle Ctrl+C
		if b == 3 {
			fmt.Println()
			term.Restore(fd, state)
			return "", fmt.Errorf("cancelled by user")
		}

		// Ignore other control characters
		if b < 32 {
			continue
		}

		// Regular character input
		if placeholderVisible {
			// Clear placeholder on first character
			fmt.Print("\033[K") // Clear to end of line
			placeholderVisible = false
		}

		// Add to buffer and display
		r := rune(b)
		inputBuffer = append(inputBuffer, r)
		fmt.Printf("%c", r)
	}

	input := string(inputBuffer)

	// Clear the prompt line
	fmt.Print("\033[1A\033[K")

	// If empty, use default
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

// PromptConfirm prompts for yes/no confirmation
func PromptConfirm(prompt string) (bool, error) {
	fmt.Printf("%s%s (y/n):%s ", ColorYellow, prompt, ColorReset)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

// PromptConfirmWithDefault prompts for yes/no confirmation with default value and placeholder
func PromptConfirmWithDefault(prompt string, defaultValue bool) (bool, error) {
	var placeholder, defaultChar string
	if defaultValue {
		placeholder = "Y/n"
		defaultChar = "y"
	} else {
		placeholder = "y/N"
		defaultChar = "n"
	}

	fmt.Printf("%s%s (%s):%s ", ColorYellow, prompt, placeholder, ColorReset)

	// Print placeholder in gray
	fmt.Printf("%s%s%s", ColorGray, defaultChar, ColorReset)
	// Move cursor back
	fmt.Printf("\033[1D")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))

	// If empty, use default
	if input == "" {
		return defaultValue, nil
	}

	return input == "y" || input == "yes", nil
}
