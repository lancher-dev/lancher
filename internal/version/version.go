package version

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
