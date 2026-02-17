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
	stdout, stderr, err := executeInitCmd(t, "-s")
	if err != nil {
		t.Fatalf("execute init -s failed: %v", err)
	}
	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
	if !strings.Contains(stdout, "scmpuff_status()") {
		t.Fatalf("expected output to contain scmpuff_status function, got: %q", stdout)
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
	tests := []struct {
		name        string
		args        []string
		wantAliasGS bool
		wantAliasGA bool
	}{
		{name: "default true", args: []string{"-s"}, wantAliasGS: true, wantAliasGA: true},
		{name: "short true", args: []string{"-as"}, wantAliasGS: true, wantAliasGA: true},
		{name: "separate true", args: []string{"-a", "-s"}, wantAliasGS: true, wantAliasGA: true},
		{name: "explicit true", args: []string{"-s", "--aliases=true"}, wantAliasGS: true, wantAliasGA: true},
		{name: "explicit false", args: []string{"-s", "--aliases=false"}, wantAliasGS: false, wantAliasGA: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := executeInitCmd(t, tt.args...)
			if err != nil {
				t.Fatalf("execute init failed: %v", err)
			}
			if stderr != "" {
				t.Errorf("unexpected stderr: %q", stderr)
			}

			gotAliasGS := strings.Contains(stdout, "alias gs='scmpuff_status'")
			if gotAliasGS != tt.wantAliasGS {
				t.Errorf("alias gs presence = %v, want %v", gotAliasGS, tt.wantAliasGS)
			}

			gotAliasGA := strings.Contains(stdout, "alias ga='git add'")
			if gotAliasGA != tt.wantAliasGA {
				t.Errorf("alias ga presence = %v, want %v", gotAliasGA, tt.wantAliasGA)
			}
		})
	}
}

func TestNewInitCmd_WrapFlagControlsOutput(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantWrap bool
	}{
		{name: "default true", args: []string{"-s"}, wantWrap: true},
		{name: "short true", args: []string{"-ws"}, wantWrap: true},
		{name: "separate true", args: []string{"-w", "-s"}, wantWrap: true},
		{name: "explicit true", args: []string{"-s", "--wrap=true"}, wantWrap: true},
		{name: "explicit false", args: []string{"-s", "--wrap=false"}, wantWrap: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := executeInitCmd(t, tt.args...)
			if err != nil {
				t.Fatalf("execute init failed: %v", err)
			}
			if stderr != "" {
				t.Errorf("unexpected stderr: %q", stderr)
			}

			gotWrap := strings.Contains(stdout, "function git()")
			if gotWrap != tt.wantWrap {
				t.Errorf("git wrapper presence = %v, want %v", gotWrap, tt.wantWrap)
			}
		})
	}
}
