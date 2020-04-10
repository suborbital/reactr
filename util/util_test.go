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
