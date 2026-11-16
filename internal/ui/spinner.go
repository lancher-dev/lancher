package ui

import (
	"fmt"
	"sync"
	"time"
)

// Spinner provides a simple ASCII spinner for long-running operations
type Spinner struct {
	message string
	frames  []string
	stop    chan bool
	wg      sync.WaitGroup
	active  bool
	mu      sync.Mutex
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"-", "\\", "|", "/"},
		stop:    make(chan bool),
		active:  false,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	if !IsTTY {
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the line
				fmt.Printf("\r\033[K")
				return
			default:
				fmt.Printf("\r%s%s%s %s", ColorYellow, s.frames[i%len(s.frames)], ColorReset, s.message)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop terminates the spinner
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	if !IsTTY {
		return
	}

	s.stop <- true
	s.wg.Wait()
}

// Success stops the spinner and shows success message
func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, message)
}
