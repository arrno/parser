package parser

func min[T int | uint](a T, b T) T {
	if a <= b {
		return a
	}
	return b
}

func max[T int | uint](a T, b T) T {
	if a >= b {
		return a
	}
	return b
}

var DefaultInstructions ParseInstructions = ParseInstructions{
	{
		OpenTag:  "<p>",
		CloseTag: "</p>",
		Attributes: map[string]any{
			"Type": "p",
		},
		Trim: true,
	},
	{
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Attributes: map[string]any{
			"Type": "b",
		},
	},
	{
		OpenTag:  "<em>",
		CloseTag: "</em>",
		Attributes: map[string]any{
			"Type": "i",
		},
	},
	{
		OpenTag:  "<code>",
		CloseTag: "</code>",
		Attributes: map[string]any{
			"Type": "c",
		},
	},
	{
		OpenTag:  "<a>",
		CloseTag: "</a>",
		Attributes: map[string]any{
			"Type": "a",
		},
		Trim: true,
	},
	{
		OpenTag:  "<h1>",
		CloseTag: "</h1>",
		Attributes: map[string]any{
			"Type": "h1",
		},
		Trim: true,
	},
	{
		OpenTag:  "<h2>",
		CloseTag: "</h2>",
		Attributes: map[string]any{
			"Type": "h2",
		},
		Trim: true,
	},
	{
		OpenTag:  "<h3>",
		CloseTag: "</h3>",
		Attributes: map[string]any{
			"Type": "h3",
		},
		Trim: true,
	},
	{
		OpenTag:  "<div>",
		CloseTag: "</div>",
		Attributes: map[string]any{
			"Type": "d",
		},
		Trim: true,
	},
	{
		OpenTag:  "<span>",
		CloseTag: "</span>",
		Attributes: map[string]any{
			"Type": "s",
		},
	},
	{
		OpenTag:  "<ol>",
		CloseTag: "</ol>",
		Attributes: map[string]any{
			"Type": "ol",
		},
		Trim: true,
	},
	{
		OpenTag:  "<li>",
		CloseTag: "</li>",
		Attributes: map[string]any{
			"Type": "li",
		},
		Trim: true,
	},
}
