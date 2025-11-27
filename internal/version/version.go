package version

import (
	"strconv"
	"strings"
)

var (
	// Version is set via ldflags at build time
	Version = "dev"
	// Commit is set via ldflags at build time
	Commit = "unknown"
)

// Get returns the current version
func Get() string {
	return Version
}

// GetFull returns version with commit info
func GetFull() string {
	if Commit != "unknown" {
		return Version + " (" + Commit[:7] + ")"
	}
	return Version
}

// Compare compares two version strings in format v0.0.1
// Returns:
//   1 if v1 > v2
//   0 if v1 == v2
//  -1 if v1 < v2
func Compare(v1, v2 string) int {
	// Remove 'v' prefix if present
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// Split by dots
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare each part
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int

		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}

		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}

	return 0
}

// IsNewer returns true if newVer is newer than currentVer
func IsNewer(currentVer, newVer string) bool {
	return Compare(newVer, currentVer) > 0
}
