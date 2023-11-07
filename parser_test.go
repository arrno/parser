package parser

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	parser := NewParser(nil)
	d := parser.DoParse(markup)
	r, err := json.MarshalIndent(d, "", "    ")
	assert.Nil(t, err)
	assert.Equal(t, expected, string(r))
}

func TestMapKeys(t *testing.T) {
	parser := NewParser(nil)
	result, parsed := parser.ParseMapKeys("::[ one: a, two: ✨blue, three: FuNtImE! ] nothing else matters!")
	expected := map[string]string{"one": "a", "two": "✨blue", "three": "FuNtImE!"}
	assert.Equal(t, parsed, 41)
	assert.True(t, reflect.DeepEqual(result, expected))
	result, parsed = parser.ParseMapKeys("hello!")
	assert.Equal(t, parsed, 0)
}

var markup string = `
<p>
Hello, ✨ world!
</p>::[injected: content, more: values]

<h1> one </h1>
<h2> two </h2>
<h3> three </h3>
<h2> two </h2>
<h1> one </h1>

<p> 
The brown <strong>fox</strong> jumps <code>over</code> the <em>lazy</em> dog. <strong>this is <code>very <em>nested</em></code>::[injected: content, more: values!!!] rem<ap>ip<n //>der </strong> 
</p>

<p>
<p> <p><p>			<p>woop</p></p></p>
</p>
</p>

<p>

   trim me

</p>

<p>
<span>

   don't trim me

</span>::[]
</p>::[injected: content]
cut out
`

var expected string = `[
    {
        "Content": "Hello, ✨ world!",
        "Type": "p",
        "injected": "content",
        "more": "values"
    },
    {
        "Content": "one",
        "Type": "h1"
    },
    {
        "Content": "two",
        "Type": "h2"
    },
    {
        "Content": "three",
        "Type": "h3"
    },
    {
        "Content": "two",
        "Type": "h2"
    },
    {
        "Content": "one",
        "Type": "h1"
    },
    {
        "Content": [
            {
                "Content": "The brown "
            },
            {
                "Content": "fox",
                "Type": "b"
            },
            {
                "Content": " jumps "
            },
            {
                "Content": "over",
                "Type": "c"
            },
            {
                "Content": " the "
            },
            {
                "Content": "lazy",
                "Type": "i"
            },
            {
                "Content": " dog. "
            },
            {
                "Content": [
                    {
                        "Content": "this is "
                    },
                    {
                        "Content": [
                            {
                                "Content": "very "
                            },
                            {
                                "Content": "nested",
                                "Type": "i"
                            }
                        ],
                        "Type": "c",
                        "injected": "content",
                        "more": "values!!!"
                    },
                    {
                        "Content": " rem\u003cap\u003eip\u003cn //\u003eder "
                    }
                ],
                "Type": "b"
            }
        ],
        "Type": "p"
    },
    {
        "Content": [
            {
                "Content": [
                    {
                        "Content": [
                            {
                                "Content": [
                                    {
                                        "Content": "woop",
                                        "Type": "p"
                                    }
                                ],
                                "Type": "p"
                            }
                        ],
                        "Type": "p"
                    }
                ],
                "Type": "p"
            }
        ],
        "Type": "p"
    },
    {
        "Content": "trim me",
        "Type": "p"
    },
    {
        "Content": [
            {
                "Content": "\n\n   don't trim me\n\n",
                "Type": "s"
            }
        ],
        "Type": "p",
        "injected": "content"
    }
]`
