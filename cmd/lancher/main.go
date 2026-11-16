package main

import (
	"fmt"
	"os"

	"github.com/lancher-dev/lancher/internal/cli"
	"github.com/lancher-dev/lancher/internal/ui"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s✗%s %v\n", ui.ColorRed, ui.ColorReset, err)
		os.Exit(1)
	}
}
