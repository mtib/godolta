package main

import (
	"fmt"
	"testing"
)

func TestNameStripper(t *testing.T) {
	example := map[string]string{
		"test.txt.delta": "test.txt",
		"delta":          "delta.dec",
		"delta.txt":      "delta.txt.dec",
	}
	for k, v := range example {
		fmt.Println(k, v)
		if v != removeDelta(k) {
			t.Fail()
		}
	}
}

func TestCreateFile(t *testing.T) {
	tryCreateFile("test/test_creation.txt")
}
