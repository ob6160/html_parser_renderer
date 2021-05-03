package main

import (
  "fmt"
  "github.com/fogleman/gg"
  "strings"
)

type DOMNode struct {
  tag string
  text string
  attributes map[string]string
  selfClosing bool
  children []*DOMNode
}

func (d DOMNode) Equal(y DOMNode) bool {
  tags := d.tag == y.tag
  text := d.text == y.text
  if !tags || !text {
    return false
  }

  nilOrEmpty := (
    len(d.children) == 0  || len(y.children) == 0) &&
    (d.children == nil || y.children == nil)

  if nilOrEmpty {
    return true
  }

  eq := true
  if len(d.children) > 0 && len(y.children) > 0 {
    for i, child := range d.children {
      other := *y.children[i]
      if !child.Equal(other) {
        eq = false
      }
    }
  } else {
    eq = false
  }

  return eq
}

func (d *DOMNode) PrintTree() {
  dc := gg.NewContext(1000, 1000)
  dc.Clear()
  dc.SetRGB(0, 0, 0)

  dc.SavePNG("out.png")
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
