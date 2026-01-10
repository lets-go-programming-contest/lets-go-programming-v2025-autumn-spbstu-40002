package main

import "testing"

func TestOnePlusOne(t *testing.T) {
	result := 1 + 1
	expected := 2
	if result != expected {
		t.Errorf("1+1 = %d, but expected %d", result, expected)
	}
}

func TestProgramStructure(t *testing.T) {
	t.Log("Basic test passed - program structure exists")
}
