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

func (d *DOMNode) ToString() string {
  if d.text != "" {
    return d.text
  }
  tag := d.tag
  attributes := d.attributes
  var attrs strings.Builder
  for key, val := range attributes {
    attrs.WriteString(fmt.Sprintf(" %s='%s'", key, val))
  }
  return fmt.Sprintf("<%s%s>", tag, attrs.String())
}

func traverse (node *DOMNode, d int) {
  indent := strings.Repeat("  ", d)
  fmt.Println(indent, node.ToString())
  for _, child := range node.children {
    traverse(child, d + 1)
  }
}

