package main

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