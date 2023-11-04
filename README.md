# Markup parser for GO!
This markup parser makes it easy to quickly define your own markup rules so you can write effortlessly and convert your content into hierarchical structured data.
## Define a set of parser rule blocks
Define a set of custom rules then use it to instantiate a parser. This example shows a set with only one instruction but you can add as many instructions as you like. If instructions is `nil`, a default subset of HTML rules are used.
```Go
instructions := ParseInstructions{
	{
		OpenTag:  "<p>",
		CloseTag: "</p>",
		Attributes: map[string]any{
			"Type": "p",
		},
		Trim: true,
	},
}

parser := NewParser(instructions)
```
## Parse your own markup language
Use your custom rules to parse text into structured data.
```Go
data := parser.DoParse("<p> Hello, World! </p>")
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

## Parse nested markup language with more tag attributes
```Go
instructions := ParseInstructions{
	{
		OpenTag:  "<p>",
		CloseTag: "</p>",
		Attributes: map[string]any{
			"Type": "p",
			"Name": "Paragraph",
		},
		Trim: true,
	},
	{
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Attributes: map[string]any{
			"Type": "b",
			"Name": "Bold",
		},
	},
	{
		OpenTag:  "<em>",
		CloseTag: "</em>",
		Attributes: map[string]any{
			"Type": "i",
			"Name": "Italic",
		},
	},
	{
		OpenTag:  "<code>",
		CloseTag: "</code>",
		Attributes: map[string]any{
			"Type": "c",
			"Name": "Code",
		},
	},
}

text := `
<p>
	The brown <strong>fox</strong> jumps <code>over</code> the <em>lazy</em> dog.
</p>
`

parser := NewParser(instructions)
data := parser.DoParse(text)
r, _ := json.MarshalIndent(data, "", "    ")
fmt.Println(string(r))

// expected:
// [
//     {
//         "Content": [
//             {
//                 "Content": "The brown "
//             },
//             {
//                 "Content": "fox",
//                 "Type": "b",
// 				   "Name": "Bold"
//             },
//             {
//                 "Content": " jumps "
//             },
//             {
//                 "Content": "over",
//                 "Type": "c",
// 				   "Name": "Code"
//             },
//             {
//                 "Content": " the "
//             },
//             {
//                 "Content": "lazy",
//                 "Type": "i",
// 				   "Name": "Italic"
//             },
//             {
//                 "Content": " dog."
//             }
//         ],
//         "Type": "p",
// 		   "Name": "Paragraph"
//     }
// ]
```
## Inject dynamic attributes
Sometimes you need to inject dynamic attributes at the individual tag level like for an `HTML` `href` or `img src`. To do this, immediately precede a closing tag with `::` followed by the key pairs within brackets and separated by commas as shown below:

```Go
instructions := ParseInstructions{
	{
		OpenTag:  "<a>",
		CloseTag: "</a>",
		Attributes: map[string]any{
			"Type": "a",
		},
		Trim: true,
	},
}
text := `<a>A link with dynamic attributes!</a>::[ href: www.example.com, dynamic: value ]`

parser := NewParser(instructions)
data := parser.DoParse(text)
r, _ := json.MarshalIndent(data, "", "    ")
fmt.Println(string(r))

// expected:
// [
//     {
//         "Content": "A link with dynamic attributes!",
//         "Type": "s",
//         "href": "www.example.com",
//         "dynamic": "value"
//     }
// ]
```