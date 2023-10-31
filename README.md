# Markup parser for GO!
This markup parser makes it easy to quickly define your own markup rules so you can write effortlessly and convert your content into hierarchical structured data.
## Define a set of parser rule blocks
Define a set of custom rules then use it to instantiate a parser. This example shows a set of instructions with one rule but you can add as many rule blocks as you like.
```Go
var instructions []*Block = []*Block{
	{
		BlockStart: "//p\n",
		BlockStop:  "\n//p",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "p",
		},
	},
}

parser := NewParser(instructions, nil)
```
## Parse your own markup language
```Go
data := parser.DoParse("//p\nHello, World!\n//p")
r, _ := json.MarshalIndent(data, "", "    ")
fmt.Println(string(r))

// expected:
// [
//     {
//         "Content": "Hello, World!",
//         "Type": "p"
//     }
// ]
```

## Parse nested markup language and inject custom attributes
Check out this more verbose example of a rule block with nested blocks. We can also use a default instruction block at any level for text not captured by tags defined in rule blocks.
```Go
var nestedInstructions []*Block = []*Block{
	{
		BlockStart: "//inlinep\n",
		BlockStop:  "\n//inlinep",
		SubBlocks: []*Block{
			{
				BlockStart: "[[",
				BlockStop:  "]]",
				SubBlocks:  nil,
				InjectValues: map[string]any{
					"Type": "code",
				},
			},
			{
				BlockStart: "**",
				BlockStop:  "**",
				SubBlocks:  nil,
				InjectValues: map[string]any{
					"Type": "b",
				},
			},
			{
				BlockStart: "__",
				BlockStop:  "__",
				SubBlocks:  nil,
				InjectValues: map[string]any{
					"Type": "i",
				},
			},
		},
		SubDefaultBlock: &Block{
			BlockStart: "<span>",
			BlockStop:  "</span>",
			SubBlocks:  nil,
			InjectValues: map[string]any{
				"Type": "span",
			},
		},
		InjectValues: map[string]any{
			"Type": "inlinep",
			"Dynamic": "Inject whatever you want!"
		},
	},
}
```

Here is an example of nested markup. We want everything inside the `inlinep` tags to be considered a paragraph. Within that paragraph, we have markup for `bold`, `italic`, and `code` text. We also want anything within the paragraph that is not contained in an explicit markup tag to be marked as a `span`. We use the rule set described above to instantiate a new parser and produce the desired structured output.
```Go
text := `
//inlinep
The brown **fox** jumps [[over]] the __lazy__ dog.
//inlinep
`

parser := NewParser(nestedInstructions, nil)
data := parser.DoParse(text)
r, _ := json.MarshalIndent(data, "", "    ")
fmt.Println(string(r))

// expected:
// [
//     {
//         "Content": [
//             {
//                 "Content": "The brown ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "fox",
//                 "Type": "b"
//             },
//             {
//                 "Content": " jumps ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "over",
//                 "Type": "code"
//             },
//             {
//                 "Content": " the ",
//                 "Type": "span"
//             },
//             {
//                 "Content": "lazy",
//                 "Type": "i"
//             },
//             {
//                 "Content": " dog.",
//                 "Type": "span"
//             }
//         ],
//         "Type": "inlinep"
//         "Dynamic": "Inject whatever you want!"
//     }
// ]
```
## Inject dynamic attributes
Sometimes you need to inject dynamic attributes at the individual tag level like for an `HTML` `href` or `img src`. To do this, immediately preceed a closing tag with `::` followed by the key pairs within brackets and separated by commas as shown below:

```Go
// assume we have defined a rule such as:
// {
// 	BlockStart: "<a>",
// 	BlockStop:  "</a>",
// 	SubBlocks:  nil,
// 	InjectValues: map[string]any{
// 		"Type": "Link",
// 	},
// }

text := `<a>A link with dynamic attributes!</a>::[ href: www.example.com, dynamic: value]`

parser := NewParser(nestedInstructions, nil)
data := parser.DoParse(text)
r, _ := json.MarshalIndent(data, "", "    ")
fmt.Println(string(r))

// expected:
// [
//     {
//         "Content": "A link with dynamic attributes!",
//         "Type": "Link",
//         "href": "www.example.com",
//         "dynamic": "value"
//     }
// ]
```