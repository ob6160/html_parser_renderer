package main

/**
 * Data structure which calculates the position of each DOM node
 * in the tree relative to one-another.
 */

type Layout struct {
	root *DOMNode
}

func NewLayout(root *DOMNode) *Layout {
	return &Layout{
	root,
	}
}