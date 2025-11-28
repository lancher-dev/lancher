package shared

import (
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	spinner := NewSpinner("Testing...")

	// Start spinner
	spinner.Start()

	// Let it run briefly
	time.Sleep(500 * time.Millisecond)

	// Stop normally
	spinner.Stop()

	// Verify we can restart
	spinner = NewSpinner("Testing again...")
	spinner.Start()
	time.Sleep(200 * time.Millisecond)
	spinner.Success("Test completed")
}

func TestSpinnerWriter(t *testing.T) {
	// Test verbose mode
	writer := NewSpinnerWriter(true)
	n, err := writer.Write([]byte("test output"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != 11 {
		t.Errorf("Expected 11 bytes written, got %d", n)
	}

	// Test buffered mode
	writer = NewSpinnerWriter(false)
	n, err = writer.Write([]byte("buffered"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != 8 {
		t.Errorf("Expected 8 bytes written, got %d", n)
	}

	output := writer.GetOutput()
	if output != "buffered" {
		t.Errorf("Expected 'buffered', got '%s'", output)
	}
}
