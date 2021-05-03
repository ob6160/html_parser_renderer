package main

import (
	"io"
	"log"
)

func IsSelfClosing(tag string) bool {
	var dict = map[string]bool{
		"area": true,
		"base": true,
		"br": true,
		"col": true,
		"embed": true,
		"hr": true,
		"img": true,
		"input": true,
		"keygen": true,
		"link": true,
		"meta": true,
		"param": true,
		"source": true,
		"track": true,
		"wbr": true,
	}
	return dict[tag]
}

func HandleIOError(err error) bool {
	if err != nil {
		if err == io.EOF {
			return true
		}
		log.Fatal(err)
		return true
	}
	return false
}