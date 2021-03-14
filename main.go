package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "os"
  _"net/http"
)

func main() {
  if len(os.Args) < 2 {
    os.Exit(1)
  }

  url := os.Args[1]
  fmt.Printf("loading: %s\n", url)

  content, err := ioutil.ReadFile(url)

  if err != nil {
    log.Fatal(err)
  }

  parser := NewParser(string(content), 0)
  parser.Parse()
}
