package ui

import (
	"fmt"
	"time"
)

// Task runs fn while showing msg: next to a spinner, or on a line of its own
// when verbose is set or stdout is not a terminal.
func Task(msg string, verbose bool, fn func() error) error {
	if verbose || !IsTTY {
		fmt.Printf("%s%s%s\n", ColorYellow, msg, ColorReset)
		return fn()
	}

	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for i := 0; ; i++ {
			fmt.Printf("\r%s%c%s %s", ColorYellow, `-\|/`[i%4], ColorReset, msg)
			select {
			case <-done:
				fmt.Print("\r\033[K")
				return
			case <-ticker.C:
			}
		}
	}()
	err := fn()
	close(done)
	<-stopped
	return err
}
