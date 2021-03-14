package main

import "fmt"

type DOMNode struct {
  tag string
  text string
  attributes map[string]string
  children []*DOMNode
}

func (d *DOMNode) PrettyPrint() {
  fmt.Println(*d)
}

