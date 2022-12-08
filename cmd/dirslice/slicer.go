package main

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Number of files to process at once.
// Smaller numbers mean more stat() calls
// Larger numbers consume more memory
const statChunkSize = 4096

// A file's location
type fileLocation struct {
	// files directory relative to a root directory
	dir *string
	// filename
	name string
}

// Recursively collect entries in a directory
func collectDirEntries(rootPath, subPath string, dirEntries []fs.DirEntry) ([]fileLocation, error) {
	var result []fileLocation

	for _, fileInfo := range dirEntries {
		if strings.HasPrefix(fileInfo.Name(), ".") {
			// Skip dot files
			continue
		} else if fileInfo.IsDir() {
			children, err := collectDir(rootPath, filepath.Join(subPath, fileInfo.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, children...)
		} else {
			result = append(result, fileLocation{dir: &subPath, name: fileInfo.Name()})
		}
	}

	return result, nil
}

// Recursively get all files in a directory, skipping .dot files
func collectDir(rootPath, subPath string) ([]fileLocation, error) {
	f, err := os.Open(filepath.Join(rootPath, subPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result []fileLocation

	for {
		dirEntries, err := f.ReadDir(statChunkSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		files, err := collectDirEntries(rootPath, subPath, dirEntries)
		if err != nil {
			return nil, err
		}
		result = append(result, files...)
	}

	return result, nil
}

// Create a target directory.
//
// A prefix of the first and last filenames in the directory is used to create
// a name in the form of FIRST-LAST
func makeDirname(startFilename, endFilename string, prefixLength int) string {
	startPrefix := getPrefix(startFilename, prefixLength)
	endPrefix := getPrefix(endFilename, prefixLength)

	return startPrefix + "-" + endPrefix
}

// Create unique directory names
func makeDirnames(files []fileLocation, sliceSize int) []string {
	numOfSlices := int(math.Ceil(float64(len(files)) / float64(sliceSize)))

	// Do not create directory names if there is just one slice
	if numOfSlices == 1 {
		return []string{}
	}

	prefixLength := int(math.Ceil(float64(numOfSlices) / 26.))

	seenDirnames := newSet[string]()
	var result []string

	for i := 0; i < numOfSlices; i++ {
		dirname := makeDirname(
			files[i*sliceSize].name,
			files[min((i+1)*sliceSize, len(files))-1].name,
			prefixLength)

		dirname = makeUnique(dirname, seenDirnames)

		seenDirnames.put(dirname)
		result = append(result, dirname)
	}

	return result
}

func sliceDir(sourceRoot, targetRoot string, sliceSize int) {
	// Collect all files from the source directory
	files, err := collectDir(sourceRoot, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	// Sort the files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	dirnames := makeDirnames(files, sliceSize)

	// If there is only one slice, just move the files into the target directory
	if len(dirnames) == 0 {
		moveFiles(sourceRoot, targetRoot, files)
	}

	// If there is more than one slice, move the files into separate directories
	for i := 0; i < len(dirnames); i++ {
		targetDir := filepath.Join(targetRoot, dirnames[i])

		os.MkdirAll(targetDir, os.ModePerm)

		start := i * sliceSize
		slice := files[start : start+min(sliceSize, len(files)-start)]

		moveFiles(sourceRoot, targetDir, slice)
	}
}

// Move files from sourceDir to targetDir.
// If a file with the same name already exists in targetDir, the file will be
// renamed
func moveFiles(sourceDir, targetDir string, files []fileLocation) {
	seenFilenames := newSet[string]()

	for _, e := range files {
		sourcePath := filepath.Join(sourceDir, *e.dir, e.name)

		targetFilename := makeUnique(e.name, seenFilenames)
		targetPath := filepath.Join(targetDir, targetFilename)

		os.Rename(sourcePath, targetPath)

		seenFilenames.put(targetFilename)
	}
}
