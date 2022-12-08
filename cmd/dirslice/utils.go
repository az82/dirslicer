package main

import (
	"fmt"
)

// Smaller of two integers
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Get the prefix of up to a given rune length
func getPrefix(str string, length int) string {
	runes := []rune(str)
	if length > len(runes) {
		length = len(runes)
	}

	return string(runes[0:length])
}

// Make a unique name by appending (INDEX)
// OtherNames must be a sorted slice of strings
func makeUnique(name string, otherNames set[string]) string {

	newFilename := name
	for count := 1; otherNames.contains(newFilename); count++ {
		newFilename = fmt.Sprintf("%s(%d)", name, count)
	}

	return newFilename
}
