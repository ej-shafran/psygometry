package main

import "testing"

func TestAddition(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("addition is hard...")
	}

	t.Log("Success!")
}
