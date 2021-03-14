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

func (d *DOMNode) printOpenTag() string {
  if d.text != "" {
    return d.text
  }
  tag := d.tag
  attributes := d.printAttributes()
  return fmt.Sprintf("<%s%s>", tag, attributes)
}

func (d *DOMNode) printCloseTag() string {
  if d.text != "" {
    return ""
  }
  tag := d.tag
  return fmt.Sprintf("</%s>", tag)
}

func (d *DOMNode) printAttributes() string {
  var attrs strings.Builder
  for key, val := range d.attributes {
    attrs.WriteString(fmt.Sprintf(" %s='%s'", key, val))
  }
  return attrs.String()
}

func traverse (node *DOMNode, d int) {
  indent := strings.Repeat("  ", d)
  fmt.Println(indent, node.printOpenTag())
  for _, child := range node.children {
    traverse(child, d + 1)
  }
  if node.text == "" {
    fmt.Println(indent, node.printCloseTag())
  }
}

