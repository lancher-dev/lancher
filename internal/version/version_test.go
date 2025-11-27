package version

import "testing"

func TestCompare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{"equal versions", "v1.0.0", "v1.0.0", 0},
		{"equal without v prefix", "1.0.0", "1.0.0", 0},
		{"v1 greater major", "v2.0.0", "v1.0.0", 1},
		{"v1 greater minor", "v1.1.0", "v1.0.0", 1},
		{"v1 greater patch", "v1.0.1", "v1.0.0", 1},
		{"v1 less major", "v1.0.0", "v2.0.0", -1},
		{"v1 less minor", "v1.0.0", "v1.1.0", -1},
		{"v1 less patch", "v1.0.0", "v1.0.1", -1},
		{"different lengths v1 longer", "v1.0.0.1", "v1.0.0", 1},
		{"different lengths v2 longer", "v1.0.0", "v1.0.0.1", -1},
		{"mixed prefixes", "v1.2.3", "1.2.3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compare(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestIsNewer(t *testing.T) {
	tests := []struct {
		name       string
		currentVer string
		newVer     string
		want       bool
	}{
		{"newer version", "v1.0.0", "v1.0.1", true},
		{"newer minor", "v1.0.0", "v1.1.0", true},
		{"newer major", "v1.0.0", "v2.0.0", true},
		{"same version", "v1.0.0", "v1.0.0", false},
		{"older version", "v1.0.1", "v1.0.0", false},
		{"without v prefix", "1.0.0", "1.0.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNewer(tt.currentVer, tt.newVer)
			if got != tt.want {
				t.Errorf("IsNewer(%q, %q) = %v, want %v", tt.currentVer, tt.newVer, got, tt.want)
			}
		})
	}
}
