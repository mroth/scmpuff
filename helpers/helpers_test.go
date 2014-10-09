package helpers

import "testing"

var testCases = []struct {
	arg      int
	expected string
}{{1, "$e1"}, {2, "$e2"}, {99, "$e99"}}

func TestIntToEnvVar(t *testing.T) {
	for _, tc := range testCases {
		actual := IntToEnvVar(tc.arg)
		if actual != tc.expected {
			t.Fatalf("IntToEnvVar(%v): expected %v, actual %v", tc.arg, tc.expected, actual)
		}
	}
}

func BenchmarkIntToEnvVar(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IntToEnvVar(5)
	}
}
