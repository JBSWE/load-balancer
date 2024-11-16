package algorithms

import "testing"

func TestRoundRobin(t *testing.T) {
	expected := true

	result := RoundRobin()

	if result != expected {
		t.Errorf("RoundRobin() = %v; want %v", result, expected)
	}
}
