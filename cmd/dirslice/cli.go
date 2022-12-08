package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Default number of files in the created directories.
const defaultSize = 1024

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s SOURCE TARGET\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr, "  SOURCE must be a directory")
	fmt.Fprintln(os.Stderr, "  TARGET must either be an empty directory or a new location")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

func checkSource(source string) {
	if source == "" {
		exitError("Source must be provided")
	}

	dStat, err := os.Stat(source)
	if os.IsNotExist(err) {
		exitError("Source does not exist")
	}

	if !dStat.IsDir() {
		exitError("Source is not a directory")
	}
}

func checkTarget(target string) {
	if target == "" {
		exitError("Target must be provided")
	}

	dStat, err := os.Stat(target)
	if os.IsExist(err) {
		if dStat.IsDir() {
			f, err := os.Open(target)
			if err != nil {
				exitError("Unable to open target directory")
			}
			defer f.Close()

			_, err = f.Readdirnames(1)
			if err == io.EOF {
				exitError("Target directory is not empty")
			}

		}
	}
}

func parseCliArgs() (string, string, int) {
	var help = flag.Bool("help", false, "Show help")
	var size = flag.Int("size", defaultSize, "Number of files per directory")

	flag.Parse()

	var args = flag.Args()
	if (len(args)) > 2 {
		exitError("Too many arguments")
	}

	if *help {
		printUsage()
		os.Exit(0)
	}

	var source = flag.CommandLine.Arg(0)
	checkSource(source)
	var target = flag.CommandLine.Arg(1)
	checkTarget(target)

	if *size < 1 {
		exitError("Size must be a positive number")
	}

	return source, target, *size
}

// Exit with an error message
func exitError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
