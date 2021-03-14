package main

type DOMNode struct {
  tag string
  text string
  attributes map[string]string
  children []*DOMNode
}



