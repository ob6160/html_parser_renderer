package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParser_BasicParse(t *testing.T) {
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

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
 * No closing tag. Expect the opening tag to be parsed as a valid dom node.
 */
func TestParser_NoCloseTag(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><b></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "b",
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
 * No closing tag, with text after opening tag.
 * Expect tag to be parsed and text stored as a child text node.
 */
func TestParser_NoCloseTagWithText(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><b>test</html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "b",
				children: []*DOMNode{
					{
						text: "test",
					},
				},
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
 * No open tag, expect the bold closing tag to be discarded in the parsed output.
 */
func TestParser_NoOpenTag(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html></b></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
 * No open tag. Expect the bold closing tag to be discarded in the parsed output.
 * Text included. Expect to be stored as a child text node.
 */
func TestParser_NoOpenTagWithText(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html>test</b></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				text: "test",
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestParser_MultipleChildren(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><div>test</div><div>test</div></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "div",
				children: []*DOMNode{
					{text: "test"},
				},
			},
			{
				tag: "div",
				children: []*DOMNode{
					{text: "test"},
				},
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
 * Expect second div to be nested inside first.
 */
func TestParser_MultipleChildrenNoCloseTags(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><div>test<div>test</html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "div",
				children: []*DOMNode{
					{text: "test"},
					{
						tag: "div",
						children: []*DOMNode{
							{text: "test"},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
 * Deep nesting
 */
func TestParser_DeepNesting(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><div><div><div><div>deeply nested</div></div></div></div></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "div",
				children: []*DOMNode{
					{
						tag: "div",
						children: []*DOMNode{
							{
								tag: "div",
								children: []*DOMNode{
									{
										tag: "div",
										children: []*DOMNode{
											{text: "deeply nested"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}


