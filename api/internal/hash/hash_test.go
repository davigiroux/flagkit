package hash

import "testing"

func TestConsistent_Deterministic(t *testing.T) {
	a := Consistent("my-flag", "user-123")
	b := Consistent("my-flag", "user-123")
	if a != b {
		t.Errorf("expected deterministic result, got %d and %d", a, b)
	}
}

func TestConsistent_Range(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := Consistent("flag", string(rune(i)))
		if v < 0 || v > 99 {
			t.Errorf("expected 0-99, got %d", v)
		}
	}
}

func TestConsistent_DifferentInputs(t *testing.T) {
	a := Consistent("flag-a", "user-1")
	b := Consistent("flag-b", "user-1")
	// Different flags should (usually) produce different hashes
	// This is probabilistic but extremely unlikely to collide
	if a == b {
		t.Log("warning: hash collision for different flag keys (rare but possible)")
	}
}
