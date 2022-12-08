package main

import (
	"os"
)

func main() {
	var source, target, size = parseCliArgs()

	// Create target directory
	err := os.MkdirAll(target, os.ModePerm)
	if err != nil {
		exitError("Unable to create target directory")
	}

	sliceDir(source, target, size)
}
