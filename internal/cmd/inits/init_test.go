package inits

import (
	"bytes"
	"strings"
	"testing"
)

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

func executeInitCmd(t *testing.T, args ...string) (stdout string, stderr string, err error) {
	t.Helper()

	cmd := NewInitCmd()
	cmd.SetArgs(args)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)

	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func TestNewInitCmd_StatusShortcutsIncluded(t *testing.T) {
	tests := []struct {
		name       string
		shell      string
		wantScript string
	}{
		{name: "bash", shell: "bash", wantScript: scriptStatusShortcuts},
		{name: "fish", shell: "fish", wantScript: scriptStatusShortcutsFish},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := executeInitCmd(t, "--shell="+tt.shell)
			if err != nil {
				t.Fatalf("execute init --shell=%s failed: %v", tt.shell, err)
			}
			if stderr != "" {
				t.Errorf("unexpected stderr: %q", stderr)
			}
			if !strings.Contains(stdout, tt.wantScript) {
				t.Errorf("expected output to contain %s status shortcuts script", tt.shell)
			}
		})
	}
}

func TestNewInitCmd_UnrecognizedShellErrors(t *testing.T) {
	_, _, err := executeInitCmd(t, "--shell=oil")
	if err == nil {
		t.Fatalf("expected error for unrecognized shell")
	}
	if !strings.Contains(err.Error(), `unrecognized shell "oil"`) {
		t.Fatalf("expected unrecognized shell error, got: %v", err)
	}
}

func TestNewInitCmd_AliasesFlagControlsOutput(t *testing.T) {
	shells := []string{"bash", "fish"}

	tests := []struct {
		name        string
		flagArgs    []string
		wantAliases bool
	}{
		{name: "default true", flagArgs: nil, wantAliases: true},
		{name: "short true", flagArgs: []string{"-a"}, wantAliases: true},
		{name: "explicit true", flagArgs: []string{"--aliases=true"}, wantAliases: true},
		{name: "explicit false", flagArgs: []string{"--aliases=false"}, wantAliases: false},
	}

	for _, shell := range shells {
		for _, tt := range tests {
			t.Run(shell+"/"+tt.name, func(t *testing.T) {
				args := append([]string{"--shell=" + shell}, tt.flagArgs...)
				stdout, stderr, err := executeInitCmd(t, args...)
				if err != nil {
					t.Fatalf("execute init failed: %v", err)
				}
				if stderr != "" {
					t.Errorf("unexpected stderr: %q", stderr)
				}

				gotAliases := strings.Contains(stdout, scriptAliases)
				if gotAliases != tt.wantAliases {
					t.Errorf("aliases script presence = %v, want %v", gotAliases, tt.wantAliases)
				}
			})
		}
	}
}

func TestNewInitCmd_WrapFlagControlsOutput(t *testing.T) {
	type shellInfo struct {
		name           string
		wantWrapScript string
	}
	shells := []shellInfo{
		{name: "bash", wantWrapScript: scriptGitWrapper},
		{name: "fish", wantWrapScript: scriptGitWrapperFish},
	}

	tests := []struct {
		name     string
		flagArgs []string
		wantWrap bool
	}{
		{name: "default true", flagArgs: nil, wantWrap: true},
		{name: "short true", flagArgs: []string{"-w"}, wantWrap: true},
		{name: "explicit true", flagArgs: []string{"--wrap=true"}, wantWrap: true},
		{name: "explicit false", flagArgs: []string{"--wrap=false"}, wantWrap: false},
	}

	for _, shell := range shells {
		for _, tt := range tests {
			t.Run(shell.name+"/"+tt.name, func(t *testing.T) {
				args := append([]string{"--shell=" + shell.name}, tt.flagArgs...)
				stdout, stderr, err := executeInitCmd(t, args...)
				if err != nil {
					t.Fatalf("execute init failed: %v", err)
				}
				if stderr != "" {
					t.Errorf("unexpected stderr: %q", stderr)
				}

				gotWrap := strings.Contains(stdout, shell.wantWrapScript)
				if gotWrap != tt.wantWrap {
					t.Errorf("git wrapper presence = %v, want %v", gotWrap, tt.wantWrap)
				}
			})
		}
	}
}
