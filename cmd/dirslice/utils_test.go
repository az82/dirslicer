package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {

	assert.Equal(t, min(1, 2), 1)
	assert.Equal(t, min(2, 1), 1)
	assert.Equal(t, min(1, 1), 1)
}

func TestGetPrefix(t *testing.T) {
	assert.Equal(t, getPrefix("ab", 0), "")
	assert.Equal(t, getPrefix("ab", 1), "a")
	assert.Equal(t, getPrefix("ab", 2), "ab")
	assert.Equal(t, getPrefix("ab", 3), "ab")
}

func TestMakeUniqueNoOther(t *testing.T) {
	assert.Equal(t, makeUnique("a", newSet[string]()), "a")
}

func TestMakeUniqueOneOther(t *testing.T) {
	// given
	name := "a"
	otherNames := newSet[string]()
	otherNames.put(name)

	// when
	newName := makeUnique(name, otherNames)

	// then
	assert.Equal(t, newName, "a(1)")
}

func TestMakeUniqueSomeOthers(t *testing.T) {
	// given
	name := "a"
	otherNames := newSet[string]()
	otherNames.put(name)
	otherNames.put(fmt.Sprintf("%s(1)", name))
	otherNames.put(fmt.Sprintf("%s(2)", name))

	// when
	newName := makeUnique(name, otherNames)

	// then
	assert.Equal(t, newName, "a(3)")
}
