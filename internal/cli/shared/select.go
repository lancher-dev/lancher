package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SelectOption represents an option in the selection menu
type SelectOption struct {
	Value string
	Label string // If empty, Value is used as label
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
			// Clear the display
			for i := 0; i <= len(options); i++ {
				fmt.Print("\r\033[K\n")
			}
			fmt.Printf("\033[%dA", len(options)+1)
			return "", fmt.Errorf("cancelled")
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

// MultiSelect prompts user to select multiple options with arrow key navigation and space bar to toggle
// Returns a slice of selected values or error if cancelled
func MultiSelect(prompt string, choices []string) ([]string, error) {
	if len(choices) == 0 {
		return nil, fmt.Errorf("no choices provided")
	}

	fd := int(os.Stdin.Fd())
	if !isTerminal(fd) {
		// Fallback to numbered multi-selection if not a terminal
		return multiSelectWithNumbers(prompt, choices)
	}

	// Set terminal to raw mode
	if err := setRawMode(fd); err != nil {
		return multiSelectWithNumbers(prompt, choices)
	}
	defer restoreTerminal(fd)

	// Hide cursor and show initial selection
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	selected := 0
	marked := make(map[int]bool)

	renderOptions := func() {
		fmt.Printf("\r\033[K%s%s%s\n", ColorCyan, prompt, ColorReset)
		fmt.Printf("\r\033[K%s(Use arrows to move, space to toggle, enter to confirm)%s\n", ColorGray, ColorReset)
		for i, choice := range choices {
			marker := " "
			if marked[i] {
				marker = "✓"
			}
			if i == selected {
				fmt.Printf("\r\033[K%s>%s [%s%s%s] %s\n", ColorGreen, ColorReset, ColorGreen, marker, ColorReset, choice)
			} else {
				fmt.Printf("\r\033[K%s•%s [%s] %s\n", ColorGray, ColorReset, marker, choice)
			}
		}
		fmt.Printf("\033[%dA", len(choices)+2)
	}

	renderOptions()

	for {
		key, err := readKey()
		if err != nil {
			return nil, err
		}

		switch key {
		case "up":
			if selected > 0 {
				selected--
				renderOptions()
			}
		case "down":
			if selected < len(choices)-1 {
				selected++
				renderOptions()
			}
		case "space":
			marked[selected] = !marked[selected]
			renderOptions()
		case "enter":
			// Clear the display
			for i := 0; i <= len(choices)+1; i++ {
				fmt.Print("\r\033[K\n")
			}
			fmt.Printf("\033[%dA", len(choices)+2)

			// Collect marked items
			var result []string
			for i := range choices {
				if marked[i] {
					result = append(result, choices[i])
				}
			}
			return result, nil
		case "ctrl+c", "q":
			// Clear the display
			for i := 0; i <= len(choices)+1; i++ {
				fmt.Print("\r\033[K\n")
			}
			fmt.Printf("\033[%dA", len(choices)+2)
			return nil, fmt.Errorf("cancelled")
		}
	}
}

// multiSelectWithNumbers provides numbered multi-selection as fallback
func multiSelectWithNumbers(prompt string, choices []string) ([]string, error) {
	fmt.Println(prompt)
	fmt.Println("(Enter numbers separated by commas, e.g., 1,3,5)")
	for i, choice := range choices {
		fmt.Printf("  %d) %s\n", i+1, choice)
	}
	fmt.Print("Enter numbers: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}, nil
	}

	// Parse comma-separated numbers
	parts := strings.Split(input, ",")
	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		choice, err := strconv.Atoi(part)
		if err != nil || choice < 1 || choice > len(choices) {
			return nil, fmt.Errorf("invalid selection: %s", part)
		}
		result = append(result, choices[choice-1])
	}

	return result, nil
}
