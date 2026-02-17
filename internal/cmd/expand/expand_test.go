package expand

import (
	"testing"
)

// Process expansion with an empty arg should be quoted so it doesnt get lost,
// special case handling that occurs in final step (to avoid escaping).
func TestProcessEmpty(t *testing.T) {
	actual := Process([]string{"a", "", "c"}, false)
	expected := "a\t''\tc"

	if actual != expected {
		t.Fatalf("ExpandEmpty: expected %v, actual %v", expected, actual)
	}
}

func TestProcessEscaping(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "spaces escaped", value: "hi mom.txt", want: "hi\\ mom.txt"},
		{name: "glob escaped", value: "refs/wip/*", want: "refs/wip/\\*"},
		{name: "paren escaped", value: "so(dumb).jpg", want: "so\\(dumb\\).jpg"},
		{name: "quote escaped", value: "\"x.txt", want: "\\\"x.txt"},
		{name: "semicolon escaped", value: "wt;af.gif", want: "wt\\;af.gif"},
		{name: "pipe escaped", value: "foo|bar", want: "foo\\|bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("e1", tt.value)
			actual := Process([]string{"1"}, false)
			if actual != tt.want {
				t.Errorf("Process([1])=%q, want %q", actual, tt.want)
			}
		})
	}
}
