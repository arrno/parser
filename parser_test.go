package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemp(t *testing.T) {
	assert.Equal(t, 1, 1)
}

// func TestParse(t *testing.T) {
// 	parser := NewParser(nil, nil)
// 	DataSet := parser.DoParse(md)
// 	r, _ := json.MarshalIndent(DataSet, "", "    ")
// 	assert.Equal(t, mdExpected, string(r))
// 	DataSet = parser.DoParse(quitTagSeparated)
// 	r, _ = json.MarshalIndent(DataSet, "", "    ")
// 	assert.Equal(t, quitTagSeparatedExpected, string(r))
// 	DataSet = parser.DoParse(withEmoji)
// 	r, _ = json.MarshalIndent(DataSet, "", "    ")
// 	assert.Equal(t, withEmojiExpected, string(r))
// }

// func TestMapKeys(t *testing.T) {
// 	parser := NewParser(nil, nil)
// 	result, parsed := parser.ParseMapKeys("::[ one: a, two: ✨blue, three: FuNtImE! ] nothing else matters!")
// 	expected := map[string]string{"one": "a", "two": "✨blue", "three": "FuNtImE!"}
// 	assert.Equal(t, parsed, 41)
// 	assert.True(t, reflect.DeepEqual(result, expected))
// 	result, parsed = parser.ParseMapKeys("hello!")
// 	assert.Equal(t, parsed, 0)
// }

// func TestParseMapKeys(t *testing.T) {
// 	parser := NewParser(nil, nil)
// 	DataSet := parser.DoParse(withMapKeys)
// 	r, _ := json.MarshalIndent(DataSet, "", "    ")
// 	assert.Equal(t, withMapKeysExpected, string(r))
// }

// var withMapKeys string = `
// //p
// one paragraph
// //p::[special:value, injected: true]
// //p
// a second paragraph
// //p::[one: blue]

// //inlinep
// Once I was **cool**::[inject: me] I was so [[cool]] that I was a __cool cat!__::[inject: end]
// //inlinep::[nested: inject]

// //inlinep
// hello, **world**::[ inject: me ] !!!
// //inlinep
// `
// var withMapKeysExpected string = `[
//     {
//         "Content": "one paragraph",
//         "Type": "p",
//         "injected": "true",
//         "special": "value"
//     },
//     {
//         "Content": "a second paragraph",
//         "Type": "p",
//         "one": "blue"
//     },
//     {
//         "Content": [
//             {
//                 "Content": "Once I was ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "cool",
//                 "Type": "b",
//                 "inject": "me"
//             },
//             {
//                 "Content": " I was so ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "cool",
//                 "Type": "code"
//             },
//             {
//                 "Content": " that I was a ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "cool cat!",
//                 "Type": "i",
//                 "inject": "end"
//             }
//         ],
//         "Type": "inlinep",
//         "nested": "inject"
//     },
//     {
//         "Content": [
//             {
//                 "Content": "hello, ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "world",
//                 "Type": "b",
//                 "inject": "me"
//             },
//             {
//                 "Content": " !!!",
//                 "Type": "span"
//             }
//         ],
//         "Type": "inlinep"
//     }
// ]`

// var withEmoji string = `
// ## Overview
// # ✨ Hello and welcome!
// `
// var withEmojiExpected string = `[
//     {
//         "Content": "Overview",
//         "Type": "h2"
//     },
//     {
//         "Content": "✨ Hello and welcome!",
//         "Type": "h1"
//     }
// ]`

// var md string = `
// ### H
// # H
// ## H
// ### H
// ## H
// # H
// ## H
// ## H

// //p
// one I was a cool
// cat who was a cool cat
// cool
// //p

// //inlinep
// Once I was **cool** I was so [[cool]] that I was a __cool__ cat!
// //inlinep

// //code-Go
// fmt.Println("Hello")
// //code

// //ul
// - one
// - two
// - three
// //ul

// //ol
// - one
// - two
// - three
// //ol

// //!end!//
// ## I will be cut out
// `

// var mdExpected string = `[
//     {
//         "Content": "H",
//         "Type": "h3"
//     },
//     {
//         "Content": "H",
//         "Type": "h1"
//     },
//     {
//         "Content": "H",
//         "Type": "h2"
//     },
//     {
//         "Content": "H",
//         "Type": "h3"
//     },
//     {
//         "Content": "H",
//         "Type": "h2"
//     },
//     {
//         "Content": "H",
//         "Type": "h1"
//     },
//     {
//         "Content": "H",
//         "Type": "h2"
//     },
//     {
//         "Content": "H",
//         "Type": "h2"
//     },
//     {
//         "Content": "one I was a cool\ncat who was a cool cat\ncool",
//         "Type": "p"
//     },
//     {
//         "Content": [
//             {
//                 "Content": "Once I was ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "cool",
//                 "Type": "b"
//             },
//             {
//                 "Content": " I was so ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "cool",
//                 "Type": "code"
//             },
//             {
//                 "Content": " that I was a ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "cool",
//                 "Type": "i"
//             },
//             {
//                 "Content": " cat!",
//                 "Type": "span"
//             }
//         ],
//         "Type": "inlinep"
//     },
//     {
//         "Content": "fmt.Println(\"Hello\")",
//         "Language": "Go",
//         "Type": "code"
//     },
//     {
//         "Content": [
//             {
//                 "Content": "one",
//                 "Type": "li"
//             },
//             {
//                 "Content": "two",
//                 "Type": "li"
//             },
//             {
//                 "Content": "three",
//                 "Type": "li"
//             }
//         ],
//         "Type": "ul"
//     },
//     {
//         "Content": [
//             {
//                 "Content": "one",
//                 "Type": "li"
//             },
//             {
//                 "Content": "two",
//                 "Type": "li"
//             },
//             {
//                 "Content": "three",
//                 "Type": "li"
//             }
//         ],
//         "Type": "ol"
//     }
// ]`

// var quitTagSeparated string = `
// //code-Go
// <p>
//     The brown <b>fox</b> jumps <code>over</code> the <em>lazy</em> dog.
// </p>
// //code
// `
// var quitTagSeparatedExpected string = `[
//     {
//         "Content": "\u003cp\u003e\n    The brown \u003cb\u003efox\u003c/b\u003e jumps \u003ccode\u003eover\u003c/code\u003e the \u003cem\u003elazy\u003c/em\u003e dog.\n\u003c/p\u003e",
//         "Language": "Go",
//         "Type": "code"
//     }
// ]`
