package inits

import "testing"

func Test_defaultShellType(t *testing.T) {
	tests := []struct {
		shellenv string
		want     string
	}{
		// supported shells at a bunch of different locations
		{"/bin/zsh", "zsh"},
		{"/usr/bin/zsh", "zsh"},
		{"/usr/local/bin/zsh", "zsh"},
		{"/bin/bash", "bash"},
		{"/usr/local/bin/fish", "fish"},

		// edge cases
		{"", "sh"},
		{"/bin/unsupported", "sh"},
	}
	for _, tt := range tests {
		t.Setenv("SHELL", tt.shellenv)
		if got := defaultShellType(); got != tt.want {
			t.Errorf("defaultShellType(%v) = %v, want %v", tt.shellenv, got, tt.want)
		}
	}
}
