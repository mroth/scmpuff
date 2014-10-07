package main

import "strconv"

func intToEnvVar(num int) string {
	return "$e" + strconv.Itoa(num)
}
