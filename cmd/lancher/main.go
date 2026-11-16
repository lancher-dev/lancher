package main

import (
	"os"

	"github.com/lancher-dev/lancher/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
