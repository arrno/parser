package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	parser := NewParser(nil, nil)
	DataSet := parser.DoParse("//p\nHello, World!\n//p")
	r, _ := json.MarshalIndent(DataSet, "", "    ")
	fmt.Println(string(r))
}
