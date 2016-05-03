package main

import (
	"bytes"
	"testing"
)

func TestInput(t *testing.T) {
	s := "test output\nline"
	stdin := bytes.NewBufferString(s)
	ch := input(stdin)

	var data []byte
	for d := range ch {
		data = append(data, d...)
	}
	if s != string(data) {
		t.Fatalf("%s(%d) != %s(%d)", s, len(s), string(data), len(data))
	}
}
