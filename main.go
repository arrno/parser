package main

import (
	"encoding/json"
	"fmt"
)

var md string = `### H
# H
## H
### H
## H
# H
## H
## H

//p
one I was a cool
man who was a cool man
cool
//p

//inlinep
Once I was **cool** I was so [[cool]] that I was a __cool__ boy!
//inlinep

//code-Go
fmt.Println("Hello")
//code

//e
## I will be cut out
`

// lists
// t := `
// \\ul
// - one
// - two
// - three
// \\ul

// \\ol
// - one
// - two
// - three
// \\ol
// `

func main() {
	parser := NewParser(nil, nil)
	DataSet := parser.DoParse(md)
	r, _ := json.MarshalIndent(DataSet, "", "    ")
	fmt.Println(string(r))
}
