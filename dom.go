package main

import (
  "fmt"
  "strings"
)

type DOMNode struct {
  tag string
  text string
  attributes map[string]string
  selfClosing bool
  children []*DOMNode
}

func (d *DOMNode) PrintTree() {
  traverse(d, 0)
}

func traverse (node *DOMNode, d int) {
  indent := strings.Repeat("  ", d)
  fmt.Println(indent, node.printOpenTag())
  for _, child := range node.children {
    traverse(child, d + 1)
  }
  if node.text == "" && !node.selfClosing {
    fmt.Println(indent, node.printCloseTag())
  }
}

/**
 * Pretty prints the node as an html open tag.
 */
func (d *DOMNode) printOpenTag() string {
  if d.text != "" {
    // If it's self closing and it has text, it's a comment.
    if d.selfClosing {
      return fmt.Sprintf("<!--%s-->", d.text)
    }
    return d.text
  }
  tag := d.tag
  attributes := d.printAttributes()
  if d.selfClosing {
    return fmt.Sprintf("<%s%s />", tag, attributes)
  }
  return fmt.Sprintf("<%s%s>", tag, attributes)
}

/**
 * Pretty prints the node as an html close tag.
 */
func (d *DOMNode) printCloseTag() string {
  if d.text != "" {
    return ""
  }
  tag := d.tag
  return fmt.Sprintf("</%s>", tag)
}

/**
 * Pretty prints the nodes' attributes.
 */
func (d *DOMNode) printAttributes() string {
  var attrs strings.Builder
  for key, val := range d.attributes {
    attrs.WriteString(fmt.Sprintf(" %s=\"%s\"", key, val))
  }
  return attrs.String()
}
