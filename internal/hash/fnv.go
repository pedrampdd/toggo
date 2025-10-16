package hash

import (
	"hash/fnv"
)

// FNVHasher implements deterministic hashing using FNV-1a algorithm
type FNVHasher struct{}

// NewFNV creates a new FNV hasher
func NewFNV() *FNVHasher {
	return &FNVHasher{}
}

// Hash returns a deterministic hash value between 0 and 99
func (h *FNVHasher) Hash(s string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(s))
	return int(hasher.Sum32() % 100)
}
