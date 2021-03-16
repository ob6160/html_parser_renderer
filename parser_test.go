package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><b></b></html>", 0, true)

	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "b",
			},
		},
	}

	if cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}


