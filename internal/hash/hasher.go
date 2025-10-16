package hash

// Hasher defines the interface for hashing strategies used in rollout
type Hasher interface {
	// Hash takes a string and returns a hash value between 0 and 99 (percentage)
	Hash(s string) int
}
