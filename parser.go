package main

import (
  "bufio"
  "fmt"
  "log"
  "strings"
)

/**
 * Takes the HTML document as an input string.
 * Instantiates a buffered string reader.
 * Consumes characters from reader as part of a recursive descent parser.
 * The grammar followed is defined in grammar.txt
 * Read more: https://en.wikipedia.org/wiki/Recursive_descent_parser
 */

type Parser struct {
  input    string
  position int
  verbose  bool
  reader   *bufio.Reader
}

func NewParser(input string, position int, verbose bool) *Parser {
  sr := strings.NewReader(input)
  var reader = bufio.NewReader(sr)
  return &Parser{
    input,
    position,
    verbose,
    reader,
  }
}

/**
 * Accept the next byte if it's any of the provided check bytes.
 */
func (p *Parser) accept(check ...byte) bool {
  var next, err = p.reader.Peek(1)
  if HandleIOError(err) {
    return false
  }
  for _, checkByte := range check {
    if checkByte == next[0] {
      _, err = p.reader.ReadByte()
      if HandleIOError(err) {
        return false
      }
      return true
    }
  }
  return false
}

/**
 * Asserts that the next len(check) bytes match the check string.
 */
func (p *Parser) assertString(check string) bool {
  var next, err = p.reader.Peek(len(check))
  if HandleIOError(err) {
    return false
  }
  return string(next) == check
}

/**
 * Peeks ahead and consumes a string if it's found.
 */
func (p *Parser) acceptString(check string) bool {
  if p.assertString(check) {
    p.reader.Discard(len(check))
    return true
  }
  return false
}

/**
 * Consumes byte provided test function returns true.
 */
func (p *Parser) acceptByteGivenTest(test func(val byte) bool) (byte, bool) {
  var next, err = p.reader.Peek(1)
  if HandleIOError(err) {
    return 0, false
  }
  valid := test(next[0])
  if valid {
    _, err = p.reader.ReadByte()
    if HandleIOError(err) {
      return 0, false
    }
    return next[0], true
  }
  return 0, false
}

/**
 * Consumes all the bytes until provided test function returns false.
 */
func (p *Parser) acceptBytesUntilTest(test func(val byte) bool) string {
  var sb strings.Builder
  var state = true
  var val byte
  for state {
    val, state = p.acceptByteGivenTest(test)
    if state {
      sb.WriteByte(val)
    }
  }
  if len(sb.String()) == 0 {
    return ""
  }
  return sb.String()
}

// Consumes nothingness: carriage returns, newlines and spaces.
func (p *Parser) consumeWhitespace() {
  state := true
  for state == true {
    state = p.accept('\r', '\n', ' ', '\t')
  }
}

func isAlphanumericOrPunctuation(check byte) bool {
  return check != '>' && check != '<' && check != '\r' && check != '\n'
}

func isAttributeSplit(check byte) bool {
  return check != '=' && check != '>' && check != ' '
}

func isSingleQuote(check byte) bool {
  return check != '\''
}

func isDoubleQuote(check byte) bool {
  return check != '"'
}

/* Actual parsing starts here */

func (p *Parser) Parse() *DOMNode {
  var rootNode = p.document()
  return rootNode
}

func (p *Parser) document() *DOMNode {
  p.consumeWhitespace()
  p.acceptString("<!DOCTYPE html>")
  var rootNode = p.node()
  return rootNode
}

/**
 * Parses a single node.
 * TODO: Omit non rendered head tags from the final DOM tree.
 */
func (p *Parser) node() *DOMNode {
  p.consumeWhitespace()

  // accept html comment if it exists.
  if c := p.comment(); c != nil {
    return c
  }

  // accept html open tag if it exists
  openTag, attributes, selfClosing := p.openTag()

  // if this isn't a comment or an html tag, it must be a piece of plaintext
  if openTag == "" {
    return p.text()
  }

  node := &DOMNode{
    tag:         openTag,
    attributes:  attributes,
    selfClosing: selfClosing,
  }

  // Self closing tags can't have children :(
  if selfClosing {
    return node
  }

  // recursively discover child nodes.
  for n := p.node(); n != nil; {
    node.children = append(node.children, n)
    n = p.node()
  }

  // accept html close tag
  // very permissive, we don't check if it exists or not.
  p.closeTag()

  return node
}

/**
 * Ruleset for accepting a single html comment.
 * TODO: Think of another way to represent comments in the DOM tree without the selfClosing hack.
 */
func (p *Parser) comment() *DOMNode {
  if !p.acceptString("<!--") {
    return nil
  }

  var sb strings.Builder
  for !p.acceptString("-->") {
    val, _ := p.reader.ReadByte()
    sb.WriteByte(val)
  }

  if p.verbose {
    fmt.Println("Comment: ", sb.String())
  }

  return &DOMNode{
    text:        sb.String(),
    selfClosing: true,
  }
}

/**
 * Ruleset for accepting any text.
 */
func (p *Parser) text() *DOMNode {
  var val = p.acceptBytesUntilTest(isAlphanumericOrPunctuation)
  if len(val) > 0 {
    node := DOMNode{
      text: val,
    }
    if p.verbose {
      fmt.Println("Consumed string: ", val)
    }
    return &node
  }
  return nil
}

/**
 * Ruleset for accepting an opening tag.
 * Return values: tag name, attributes, self-closing
 */
func (p *Parser) openTag() (string, map[string]string, bool) {
  // if it's a close tag, bail out.
  if p.assertString("</") {
    return "", nil, false
  }

  if !p.accept('<') {
    return "", nil, false
  }

  var tagName = p.tagName()

  if p.verbose {
    fmt.Println("Open tag: ", tagName)
  }

  p.consumeWhitespace()

  // Exit early if we've reached the end of the tag.
  // Return true if it's a self closing tag.
  if p.accept('>') {
    selfClosing := IsSelfClosing(tagName)
    return tagName, nil, selfClosing
  } else if p.acceptString("/>") {
    return tagName, nil, true
  }

  // Tag isn't over yet, cover any attributes we can find.
  var attributes = make(map[string]string)
  var attrName, attrValue = p.attribute()
  for attrName != "" {
    p.consumeWhitespace()
    attributes[attrName] = attrValue
    // Quit the loop when we find the end of the tag.
    if p.accept('>') {
      selfClosing := IsSelfClosing(tagName)
      return tagName, attributes, selfClosing
    } else if p.acceptString("/>") {
      return tagName, attributes, true
    }
    attrName, attrValue = p.attribute()
  }

  // We never found the end of the tag :(
  return "", nil, false
}

/**
 * Ruleset for accepting a closing tag.
 */
func (p *Parser) closeTag() string {
  if !p.acceptString("</") {
    return ""
  }

  var tagName = p.tagName()

  if p.verbose {
    fmt.Println("Close tag: ", tagName)
  }

  if p.accept('>') {
    return tagName
  }

  return ""
}

/**
 * Ruleset for accepting a tag name.
 */
func (p *Parser) tagName() string {
  var tagName = p.acceptBytesUntilTest(func(val byte) bool {
    return val != '>' && val != ' ' && val != '/'
  })

  // Now we've determined the tag name, consume any remaining whitespace.
  p.consumeWhitespace()
  return tagName
}

/**
 * Ruleset for accepting a single html attribute.
 */
func (p *Parser) attribute() (string, string) {
  var attributeName = p.acceptBytesUntilTest(isAttributeSplit)

  // We have an attribute without a value
  if !p.accept('=') {
    if p.verbose {
      log.Println("Attribute: ", "(", attributeName, ")")
    }
    return attributeName, ""
  }

  p.consumeWhitespace()

  var attributeValue = ""
  if p.accept('"') {
    attributeValue = p.acceptBytesUntilTest(isDoubleQuote)
  } else if p.accept('\'') {
    attributeValue = p.acceptBytesUntilTest(isSingleQuote)
  } else {
    return "", ""
  }

  if !p.accept('"', '\'') {
    return "", ""
  }

  if p.verbose {
    log.Println("Attribute: ", "(", attributeName, "=", attributeValue, ")")
  }
  return attributeName, attributeValue
}
