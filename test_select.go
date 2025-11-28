package main

import (
	"fmt"
	"github.com/Kasui92/lancher/internal/cli/shared"
)

func main() {
	fmt.Println("Testing colored select prompt...")
	fmt.Println()

	options := []string{
		"react-template",
		"vue-template",
		"next-template",
		"express-api",
	}

	selected, err := shared.Select("Choose a template:", options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\n%sYou selected:%s %s\n", shared.ColorGreen, shared.ColorReset, selected)
}
