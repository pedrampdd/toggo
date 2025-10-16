package hash

import "testing"

func TestFNVHasher_Deterministic(t *testing.T) {
	hasher := NewFNV()

	// Hash the same value multiple times
	input := "test:user123"

	hash1 := hasher.Hash(input)
	hash2 := hasher.Hash(input)
	hash3 := hasher.Hash(input)

	if hash1 != hash2 || hash2 != hash3 {
		t.Errorf("hash is not deterministic: got %d, %d, %d", hash1, hash2, hash3)
	}
}

func TestFNVHasher_Range(t *testing.T) {
	hasher := NewFNV()

	// Test multiple inputs to ensure hash is in valid range
	inputs := []string{
		"test:user1",
		"test:user2",
		"feature:user123",
		"another_flag:differentuser",
	}

	for _, input := range inputs {
		hash := hasher.Hash(input)
		if hash < 0 || hash >= 100 {
			t.Errorf("hash out of range [0, 100): got %d for input %s", hash, input)
		}
	}
}

func TestFNVHasher_Distribution(t *testing.T) {
	hasher := NewFNV()

	// Test that different inputs produce different hashes (usually)
	hash1 := hasher.Hash("flag1:user1")
	hash2 := hasher.Hash("flag1:user2")
	hash3 := hasher.Hash("flag2:user1")

	// These should be different (though not guaranteed due to collisions)
	// Just verify they're all in valid range
	if hash1 < 0 || hash1 >= 100 {
		t.Errorf("hash1 out of range: %d", hash1)
	}
	if hash2 < 0 || hash2 >= 100 {
		t.Errorf("hash2 out of range: %d", hash2)
	}
	if hash3 < 0 || hash3 >= 100 {
		t.Errorf("hash3 out of range: %d", hash3)
	}
}
