package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	md := "//p\nHello, World!\n//p"
	parser := NewParser(nil, nil)
	DataSet := parser.DoParse(md)
	r, _ := json.MarshalIndent(DataSet, "", "    ")
	fmt.Println(string(r))
}
