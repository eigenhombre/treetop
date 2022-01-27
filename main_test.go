package main

import (
	"testing"
)

func TestTopOfPath(t *testing.T) {
	var tests = []struct {
		targetPath, path, expected string
	}{
		{".", "a", "a"},
		{".", "a/b", "a"},
		{".", "a/b/c", "a"},
		{"a", "a/b/c", "b"},
	}
	for _, test := range tests {
		got := topOfPath(test.targetPath, test.path)
		if got != test.expected {
			t.Errorf("topOfPath(%q, %q) = %q, want %q", test.targetPath, test.path, got, test.expected)
		}
	}
}

func TestCommafiedInt(t *testing.T) {
	var tests = []struct {
		in            int
		got, expected string
	}{
		{0, "0", "0"},
	}
	for _, test := range tests {
		got := commafiedInt(test.in)
		if got != test.expected {
			t.Errorf("commafiedInt(%d) = %q, want %q", test.in, got, test.expected)
		}
	}
}
