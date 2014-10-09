package helpers

import "strconv"

// IntToEnvVar takes a numeric digit and returns the scmpuffenvironment variable
// for that case.
//
// The default is $e[N], that may change in the future or become configurable.
func IntToEnvVar(num int) string {
	return "$e" + strconv.Itoa(num)
}
