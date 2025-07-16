package expand

import (
	"testing"
)

// Process expansion with an empty arg should be quoted so it doesnt get lost,
// special case handling that occurs in final step (to avoid escaping).
func TestProcessEmpty(t *testing.T) {
	actual := Process([]string{"a", "", "c"})
	expected := "a\t''\tc"

	if actual != expected {
		t.Fatalf("ExpandEmpty: expected %v, actual %v", expected, actual)
	}
}
