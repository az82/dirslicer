package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetEmpty(t *testing.T) {
	assert.True(t, !newSet[string]().contains(""), "Set is empty")
}

func TestSet(t *testing.T) {
	// given
	set := newSet[string]()
	set.put("a")

	// expect
	assert.True(t, set.contains("a"), "Test value is present")
	assert.True(t, !set.contains("b"), "Set does not contain other values")
}
