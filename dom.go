package main

import (
  "fmt"
  "strings"
)

type DOMNode struct {
  tag string
  text string
  attributes map[string]string
  children []*DOMNode
}

func (d *DOMNode) PrettyPrint() {
  traverse(d, 0)
}


func traverse (node *DOMNode, d int) {
  indent := strings.Repeat("\t", d)
  fmt.Println(indent, node)
  for _, child := range node.children {
    traverse(child, d + 1)
  }
}

