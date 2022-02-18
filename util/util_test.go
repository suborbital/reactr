package util

import (
	"testing"
)

func TestGenerateResultID(t *testing.T) {
	id := GenerateResultID()

	if len(id) != 24 {
		t.Errorf("id has length %d, expected 24", len(id))
	}
}

func TestContainsInt(t *testing.T) {
	container := []int{1, 2, 3, 4}

	if !ContainsInt(3, container) {
		t.Errorf("expected value not found in container")
	}

	if ContainsInt(5, container) {
		t.Errorf("should not have found value in container")
	}
}
