package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectDirEmpty(t *testing.T) {
	// given
	rootDir, err := os.MkdirTemp("", "collect_test_*")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(rootDir)

	// when
	fileLocations, err := collectDir(rootDir, ".")

	// then
	assert.NoError(t, err)
	assert.ElementsMatch(t, fileLocations, []fileLocation{})
}

func TestCollectDirDepth1(t *testing.T) {
	// given
	rootDir, err := os.MkdirTemp("", "collect_test_*")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(rootDir)

	thisDir := "."
	var expectedFileLocations []fileLocation
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("file-%d", i)
		touch(t, filepath.Join(rootDir, filename))
		expectedFileLocations = append(expectedFileLocations, fileLocation{dir: &thisDir, name: filename})
	}

	// when
	fileLocations, err := collectDir(rootDir, thisDir)

	// then
	assert.NoError(t, err)
	assert.ElementsMatch(t, fileLocations, expectedFileLocations)
}

func TestCollectDirDepth2(t *testing.T) {
	// given
	rootDir, err := os.MkdirTemp("", "collect_test_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(rootDir)

	thisDir := "."
	var expectedFileLocations []fileLocation
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("file-%d", i)
		touch(t, filepath.Join(rootDir, filename))
		expectedFileLocations = append(expectedFileLocations, fileLocation{dir: &thisDir, name: filename})
	}

	subDir := "dir-1"
	subPath := filepath.Join(rootDir, "dir-1")
	err = os.Mkdir(subPath, fs.ModePerm)
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("file-%d", i)
		touch(t, filepath.Join(subPath, filename))
		expectedFileLocations = append(expectedFileLocations, fileLocation{dir: &subDir, name: filename})
	}

	// when
	fileLocations, err := collectDir(rootDir, thisDir)

	// then
	assert.NoError(t, err)
	assert.ElementsMatch(t, fileLocations, expectedFileLocations)
}

func touch(t *testing.T, path string) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	file.Close()
}

func TestMakeDirnames(t *testing.T) {
	// given
	rootDir := ""

	var files []fileLocation
	for i := 0; i < 30; i++ {
		files = append(files, fileLocation{dir: &rootDir, name: fmt.Sprintf("%02d", i)})
	}

	// when
	dirnames := makeDirnames(files, 10)

	// then
	assert.ElementsMatch(t, dirnames, []string{"0-0", "1-1", "2-2"})
}

func TestSliceDir(t *testing.T) {
	// given
	sourceDir, err := os.MkdirTemp("", "fdiv_test_source_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(sourceDir)

	touch(t, filepath.Join(sourceDir, "a"))
	touch(t, filepath.Join(sourceDir, "b"))
	touch(t, filepath.Join(sourceDir, "c"))
	touch(t, filepath.Join(sourceDir, "d"))
	touch(t, filepath.Join(sourceDir, "e"))

	targetDir, err := os.MkdirTemp("", "fdiv_test_target_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(targetDir)

	// when
	sliceDir(sourceDir, targetDir, 3)

	// then
	assert.FileExists(t, filepath.Join(targetDir, "a-c", "a"))
	assert.FileExists(t, filepath.Join(targetDir, "a-c", "b"))
	assert.FileExists(t, filepath.Join(targetDir, "a-c", "c"))
	assert.FileExists(t, filepath.Join(targetDir, "d-e", "d"))
	assert.FileExists(t, filepath.Join(targetDir, "d-e", "e"))
}

func TestSliceDirOne(t *testing.T) {
	// given
	sourceDir, err := os.MkdirTemp("", "fdiv_test_source_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(sourceDir)

	touch(t, filepath.Join(sourceDir, "a"))

	targetDir, err := os.MkdirTemp("", "fdiv_test_target_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(targetDir)

	// when
	sliceDir(sourceDir, targetDir, 3)

	// then
	assert.FileExists(t, filepath.Join(targetDir, "a"))
}

func TestSliceDirEmpty(t *testing.T) {
	// given
	sourceDir, err := os.MkdirTemp("", "fdiv_test_source_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(sourceDir)

	targetDir, err := os.MkdirTemp("", "fdiv_test_target_*")
	if err != nil {
		t.Fatalf("Error creating directory: %s", err)
	}
	defer os.RemoveAll(targetDir)

	// when
	sliceDir(sourceDir, targetDir, 3)

	// then
	assertEmpty(t, targetDir)
}

func assertEmpty(t *testing.T, dirPath string) {
	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		t.Fatalf("Error opening directory: %s", err)
	}
	defer dir.Close()

	// Read the directory entries
	entries, err := dir.Readdir(0)
	if err != nil {
		t.Fatalf("Error reading directory: %s", err)
	}

	// Check if the directory is empty
	assert.Empty(t, entries, "Directory %s is not empty", dirPath)
}
