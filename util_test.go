package main

import (
	"fmt"
	"testing"
)

func TestReadFile(t *testing.T) {
	lines := ReadFile("peers.txt")
	fmt.Println(lines)
}
