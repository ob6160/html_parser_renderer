package main

import (
  "bufio"
  "fmt"
  "log"
  "strings"
  "unicode"
)

type Parser struct {
  input string
  position int
  reader *bufio.Reader
}

func NewParser(input string, position int) *Parser {
  sr := strings.NewReader(input)
  var reader = bufio.NewReader(sr)
  return &Parser{
    input,
    position,
    reader,
  }
}

/**
 * Accept the next byte if it's any of the provided check bytes.
 */
func (p *Parser) accept(check ...byte) bool {
  var next, err = p.reader.Peek(1)
  if err != nil {
    log.Fatal(err)
  }
  for _, checkByte := range check {
    if checkByte == next[0] {
      _, err = p.reader.ReadByte()
      if err != nil {
        log.Fatal(err)
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
  if err != nil {
    log.Fatal(err)
  }
  return string(next) == check
}

/**
 * Peeks ahead and consumes a string if it's found.
 */
func (p *Parser) acceptString(check string) bool {
  if p.assertString(check) {
    var readUntil = check[len(check) - 1]
    p.reader.ReadBytes(readUntil)
    return true
  }
  return false
}

/**
 * Consumes byte provided test function returns true.
 */
func (p *Parser) acceptByteGivenTest(test func(val byte) bool) (byte, bool) {
  var next, err = p.reader.Peek(1)
  if err != nil {
    log.Fatal(err)
  }
  valid := test(next[0])
  if valid {
    _, err = p.reader.ReadByte()
    if err != nil {
      log.Fatal(err)
    }
    return next[0], true
  }
  return 0, false
}

/**
 * Consumes all the bytes until provided test function returns true.
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
  return unicode.IsLetter(rune(check)) ||
    unicode.IsNumber(rune(check)) ||
    unicode.IsPunct(rune(check)) ||
    check == ' ' ||
    check == '\r' ||
    check == '\n'
}

func isAttributeSplit(check byte) bool {
  return check != '=' && check != '>' && check != ' '
}

func isQuote(check byte) bool {
  return check != '"' && check != '\''
}

/* Actual parsing starts here */

func (p *Parser) Parse() *DOMNode {
  var rootNode = p.document()
  return rootNode
}

func (p *Parser) document() *DOMNode {
  p.acceptString("<!DOCTYPE html>")
  var rootNode = p.node()
  return rootNode
}

/**
 * Parses a single node.
 */
func (p *Parser) node() *DOMNode {
  p.consumeWhitespace()

  openTag, attributes := p.openTag()
  if openTag == "" {
    return nil
  }

  var children []*DOMNode

  var n *DOMNode
  for {
    n = p.node()
    if n == nil {
      n = p.text()
      if n == nil {
        break
      }
    }
    children = append(children, n)
  }
  

  var closeTag = p.closeTag()
  if closeTag == "" {
    return nil
  }

  n = &DOMNode{
    children: children,
    tag: openTag,
    attributes: attributes,
  }

  return n
}

func (p *Parser) text() *DOMNode {
   var val = p.acceptBytesUntilTest(isAlphanumericOrPunctuation)
   if len(val) > 0 {
     node := DOMNode{
       text: val,
     }
     fmt.Println("Consumed string: ", val)
     return &node
   }
   return nil
}

func (p *Parser) openTag() (string, map[string]string) {
  // if it's a close tag, bail out.
  if p.assertString("</") {
    return "", nil
  }

  if !p.accept('<') {
    return "", nil
  }

  var tagName = p.tagName()
  fmt.Println("Open tag: ", tagName)

  // Exit early if we've reached the end of the tag.
  if p.accept('>') {
    return tagName, nil
  }

  // Tag isn't over yet, cover any attributes we can find.
  var attributes = make(map[string]string)
  var attrName, attrValue = p.attribute()
  for attrName != "" {
    attributes[attrName] = attrValue
    p.consumeWhitespace()
    // Quit the loop when we find the end of the tag.
    if p.accept('>') {
      return tagName, attributes
    }
  }

  // We never found the end of the tag :(
  return "", nil
}

func (p *Parser) closeTag() string {
  if !p.acceptString("</") {
    return ""
  }
  
  var tagName = p.tagName()

  fmt.Println("Close tag: ", tagName)

  if p.accept('>') {
    return tagName
  }

  return ""
}

func (p *Parser) tagName() string {
  var tagName = p.acceptBytesUntilTest(func(val byte) bool {
    return !(val == '>' || val == ' ')
  })

  // Now we've determined the tag name, consume any remaining whitespace.
  p.consumeWhitespace()
  return tagName
}

func (p *Parser) attribute() (string, string) {
  var attributeName = p.acceptBytesUntilTest(isAttributeSplit)

  // We have an attribute without a value
  if !p.accept('=') {
    log.Println("Attribute: ", "(", attributeName, ")")
    return attributeName, ""
  }
  
  p.consumeWhitespace()
  
  if !p.accept('"', '\'') {
    return "", ""
  }

  var attributeValue = p.acceptBytesUntilTest(isQuote)

  if !p.accept('"', '\'') {
    return "", ""
  }

  log.Println("Attribute: ", "(", attributeName, "=", attributeValue, ")")
  return attributeName, attributeValue
}
