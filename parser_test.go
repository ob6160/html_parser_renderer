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

/**
 * Self closing tags
 */
func TestParser_SelfClosingTags(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html>test<br/>self<br/>closing<hr/>tags</html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				text: "test",
			},
			{
				tag: "br",
			},
			{
				text: "self",
			},
			{
				tag: "br",
			},
			{
				text: "closing",
			},
			{
				tag: "hr",
			},
			{
				text: "tags",
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}

/**
* Self closing tags
*/
func TestParser_Comments(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><!--another test comment--><b>testing comments<!--<h1>commented</h1>--></b><!-- a test comment --></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				selfClosing: true,
				text: "another test comment",
			},
			{
				tag: "b",
				children: []*DOMNode{
					{
						text: "testing comments",
					},
					{
						selfClosing: true,
						text: "<h1>commented</h1>",
					},
				},
			},
			{
				selfClosing: true,
				text: " a test comment ",
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}


/**
* Self closing tags
 */
func TestParser_Attributes(t *testing.T) {
	parser := NewParser("<!DOCTYPE html><html><head><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\" /></head><div class=\"test\" id='main' disabled></div></html>", 0, true)
	var got = *parser.Parse()
	want := DOMNode{
		tag: "html",
		children: []*DOMNode{
			{
				tag: "head",
				children: []*DOMNode{
					{
						tag: "meta",
						attributes: map[string]string{
							"name": "viewport",
							"content": "width=device-width,initial-scale=1",
						},
					},
				},
			},
			{
				tag: "div",
				attributes: map[string]string{
					"class": "test",
					"id": "main",
					"disabled": "",
				},
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Errorf("got %v want %v", got, want)
	}
}


