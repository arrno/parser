package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	parser := NewParser(nil, nil)
	DataSet := parser.DoParse(md)
	r, _ := json.MarshalIndent(DataSet, "", "    ")
	assert.Equal(t, expectedJson, string(r))
	DataSet = parser.DoParse(txtCode)
	r, _ = json.MarshalIndent(DataSet, "", "    ")
	assert.Equal(t, expectedCode, string(r))
}

var md string = `
### H
# H
## H
### H
## H
# H
## H
## H

//p
one I was a cool
cat who was a cool cat
cool
//p

//inlinep
Once I was **cool** I was so [[cool]] that I was a __cool__ cat!
//inlinep

//code-Go
fmt.Println("Hello")
//code

//ul
- one
- two
- three
//ul

//ol
- one
- two
- three
//ol

//!end!//
## I will be cut out
`

var expectedJson string = `[
    {
        "Content": "H",
        "Type": "h3"
    },
    {
        "Content": "H",
        "Type": "h1"
    },
    {
        "Content": "H",
        "Type": "h2"
    },
    {
        "Content": "H",
        "Type": "h3"
    },
    {
        "Content": "H",
        "Type": "h2"
    },
    {
        "Content": "H",
        "Type": "h1"
    },
    {
        "Content": "H",
        "Type": "h2"
    },
    {
        "Content": "H",
        "Type": "h2"
    },
    {
        "Content": "one I was a cool\ncat who was a cool cat\ncool",
        "Type": "p"
    },
    {
        "Content": [
            {
                "Content": "Once I was ",
                "Type": "span"
            },
            {
                "Content": "cool",
                "Type": "b"
            },
            {
                "Content": " I was so ",
                "Type": "span"
            },
            {
                "Content": "cool",
                "Type": "code"
            },
            {
                "Content": " that I was a ",
                "Type": "span"
            },
            {
                "Content": "cool",
                "Type": "i"
            },
            {
                "Content": " cat!",
                "Type": "span"
            }
        ],
        "Type": "inlinep"
    },
    {
        "Content": "fmt.Println(\"Hello\")",
        "Language": "Go",
        "Type": "code"
    },
    {
        "Content": [
            {
                "Content": "one",
                "Type": "li"
            },
            {
                "Content": "two",
                "Type": "li"
            },
            {
                "Content": "three",
                "Type": "li"
            }
        ],
        "Type": "ul"
    },
    {
        "Content": [
            {
                "Content": "one",
                "Type": "li"
            },
            {
                "Content": "two",
                "Type": "li"
            },
            {
                "Content": "three",
                "Type": "li"
            }
        ],
        "Type": "ol"
    }
]`

var txtCode string = `
//code-Go
<p>
    The brown <b>fox</b> jumps <code>over</code> the <em>lazy</em> dog.
</p>
//code
`
var expectedCode string = `[
    {
        "Content": "\u003cp\u003e\n    The brown \u003cb\u003efox\u003c/b\u003e jumps \u003ccode\u003eover\u003c/code\u003e the \u003cem\u003elazy\u003c/em\u003e dog.\n\u003c/p\u003e",
        "Language": "Go",
        "Type": "code"
    }
]`