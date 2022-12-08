package main

// set data type
type set[T comparable] struct {
	delegate map[T]bool
}

// Create a new set
func newSet[T comparable]() set[T] {
	return set[T]{delegate: make(map[T]bool)}
}

// put an item into the set
func (s set[T]) put(v T) {
	s.delegate[v] = true
}

// Check whether the set contains an item
func (s set[T]) contains(v T) bool {
	_, contains := s.delegate[v]
	return contains
}
