package shared

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Select asks the user to pick one of choices, with arrow keys on a terminal
// and by number otherwise.
func Select(prompt string, choices []string) (string, error) {
	if len(choices) == 0 {
		return "", errors.New("no choices provided")
	}
	if stdinIsTerminal() {
		if restore, err := rawMode(); err == nil {
			defer restore()
			return selectWithArrows(prompt, choices)
		}
	}
	return selectWithNumbers(prompt, choices)
}

// MultiSelect asks the user to pick any number of choices, toggled with space
// on a terminal and given as comma-separated numbers otherwise.
func MultiSelect(prompt string, choices []string) ([]string, error) {
	if len(choices) == 0 {
		return nil, errors.New("no choices provided")
	}
	if stdinIsTerminal() {
		if restore, err := rawMode(); err == nil {
			defer restore()
			return multiSelectWithArrows(prompt, choices)
		}
	}
	return multiSelectWithNumbers(prompt, choices)
}

// menu draws a list in raw mode and redraws it in place.
type menu struct {
	lines int
}

func (m *menu) draw(lines []string) {
	for _, line := range lines {
		fmt.Printf("\r\033[K%s\n", line)
	}
	m.lines = len(lines)
	fmt.Printf("\033[%dA", m.lines)
}

func (m *menu) clear() {
	for range m.lines {
		fmt.Print("\r\033[K\n")
	}
	fmt.Printf("\033[%dA", m.lines)
}

func selectWithArrows(prompt string, choices []string) (string, error) {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	var m menu
	selected := 0
	render := func() {
		lines := []string{ColorCyan + prompt + ColorReset}
		for i, choice := range choices {
			if i == selected {
				lines = append(lines, fmt.Sprintf("%s>%s %s", ColorGreen, ColorReset, choice))
			} else {
				lines = append(lines, fmt.Sprintf("%s•%s %s", ColorGray, ColorReset, choice))
			}
		}
		m.draw(lines)
	}
	render()

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
			if selected < len(choices)-1 {
				selected++
				render()
			}
		case "enter":
			m.clear()
			return choices[selected], nil
		case "ctrl+c", "q":
			m.clear()
			return "", ErrCancelled
		}
	}
}

func multiSelectWithArrows(prompt string, choices []string) ([]string, error) {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	var m menu
	selected := 0
	marked := make([]bool, len(choices))
	render := func() {
		lines := []string{
			ColorCyan + prompt + ColorReset,
			ColorGray + "(Use arrows to move, space to toggle, enter to confirm)" + ColorReset,
		}
		for i, choice := range choices {
			mark := " "
			if marked[i] {
				mark = "✓"
			}
			if i == selected {
				lines = append(lines, fmt.Sprintf("%s>%s [%s%s%s] %s", ColorGreen, ColorReset, ColorGreen, mark, ColorReset, choice))
			} else {
				lines = append(lines, fmt.Sprintf("%s•%s [%s] %s", ColorGray, ColorReset, mark, choice))
			}
		}
		m.draw(lines)
	}
	render()

	for {
		key, err := readKey()
		if err != nil {
			return nil, err
		}
		switch key {
		case "up":
			if selected > 0 {
				selected--
				render()
			}
		case "down":
			if selected < len(choices)-1 {
				selected++
				render()
			}
		case " ":
			marked[selected] = !marked[selected]
			render()
		case "enter":
			m.clear()
			var result []string
			for i, choice := range choices {
				if marked[i] {
					result = append(result, choice)
				}
			}
			return result, nil
		case "ctrl+c", "q":
			m.clear()
			return nil, ErrCancelled
		}
	}
}

func selectWithNumbers(prompt string, choices []string) (string, error) {
	fmt.Println(prompt)
	for i, choice := range choices {
		fmt.Printf("  %d) %s\n", i+1, choice)
	}
	fmt.Print("Enter number: ")

	input, err := readLine("")
	if err != nil {
		return "", err
	}
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(choices) {
		return "", fmt.Errorf("invalid selection: %s", input)
	}
	return choices[n-1], nil
}

func multiSelectWithNumbers(prompt string, choices []string) ([]string, error) {
	fmt.Println(prompt)
	fmt.Println("(Enter numbers separated by commas, e.g., 1,3,5)")
	for i, choice := range choices {
		fmt.Printf("  %d) %s\n", i+1, choice)
	}
	fmt.Print("Enter numbers: ")

	input, err := readLine("")
	if err != nil || input == "" {
		return nil, err
	}
	var result []string
	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 || n > len(choices) {
			return nil, fmt.Errorf("invalid selection: %s", part)
		}
		result = append(result, choices[n-1])
	}
	return result, nil
}
