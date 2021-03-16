package main

import (
	"go-cmp"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html></html>", 0, true)

	got := parser.Parse()
	want := &DOMNode{
		tag: "html",
	}
	if cmp.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}

}


