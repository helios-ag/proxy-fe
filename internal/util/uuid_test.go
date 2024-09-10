package util

import (
	"github.com/google/uuid"
	"testing"
)

func TestUuid(t *testing.T) {
	uuidStr := Uuid()
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("Expected valid UUID, got %s", uuidStr)
	}
}
