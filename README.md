# Markdown parser for GO!
This markdown parser makes it easy to quickly define your own markup rules so you can write effortlessly and convert your content into hierarchical structured data.
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
var instructions []*Block = []*Block{
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

parser := NewParser(instructions, nil)
```