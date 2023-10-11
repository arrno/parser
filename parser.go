package parser

import (
	"reflect"
	"strings"
)

// MarkupParser is something that is able to parse markup text into structured data.
type MarkupParser interface {
	DoParse(markup string) []map[string]any
}

// Block is a single markup parsing instruction. A block can have nested sub blocks.
type Block struct {
	BlockStart        string
	BlockStop         string
	Matching          bool
	MatchIndex        int
	SubBlocks         []*Block
	SubDefaultBlock   *Block
	InjectValues      map[string]any
	ContentStartIndex int
}

// Parser implements MarkupParser interface
type Parser struct {
	instructions []*Block
	defaultBlock *Block
}

// NewParser creates a new Parser and returns a pointer to it.
func NewParser(instructions []*Block, defaultBlock *Block) *Parser {
	if instructions == nil {
		instructions = DefaultInstructions
	}
	p := Parser{
		instructions: instructions,
		defaultBlock: defaultBlock,
	}
	return &p
}

// DoParse will parse the input markup text into structured data.
func (p *Parser) DoParse(markup string) []map[string]any {
	data, _ := p.handleParse(p.instructions, markup, "//!end!//", p.defaultBlock)
	return data
}

// handleParse is a recursive function that executes the parsing argorithm.
func (p *Parser) handleParse(instructions []*Block, markup, fullQuit string, defaultBlock *Block) ([]map[string]any, int) {

	dataSet := []map[string]any{}
	fullQuitMatch := 0
	var activeBlock *Block
	var index int

	blockIsActive := func(block *Block) bool {
		return reflect.DeepEqual(activeBlock, block)
	}

	nonMatchStart := 0
	var candidateBlock *Block = nil
	killSwitch := false
	skipTo := 0

	for _, char := range markup {
		index++
		i := index - 1

		// handle for recursion
		if []rune(fullQuit)[fullQuitMatch] == char {
			fullQuitMatch++
		} else {
			fullQuitMatch = 0
		}
		if len([]rune(fullQuit)) == fullQuitMatch {
			killSwitch = true
			break
		}

		for _, block := range instructions {

			if skipTo > i {
				continue
			}

			// when an active block is selected
			if blockIsActive(block) && block.SubBlocks != nil {
				// recursion case
				var content []map[string]any
				var skipBy int
				content, skipBy = p.handleParse(block.SubBlocks, string([]rune(markup)[i:len([]rune(markup))]), block.BlockStop, block.SubDefaultBlock)
				skipTo = i + skipBy
				data := map[string]any{
					"Content": content,
				}
				for k, v := range block.InjectValues {
					data[k] = v
				}
				dataSet = append(dataSet, data)
				activeBlock = nil
				p.resetBlocks(instructions, nil)
				nonMatchStart = i + 1
				continue
			}
			if activeBlock != nil && !blockIsActive(block) {
				continue
			} else if activeBlock != nil && block.ContentStartIndex == 0 {
				block.ContentStartIndex = i
			}
			if blockIsActive(block) && len([]rune(block.BlockStop)) == (block.MatchIndex+1) && []rune(block.BlockStop)[block.MatchIndex] == char {
				// append to DataSet
				data := map[string]any{
					"Content": string([]rune(markup)[block.ContentStartIndex:(i - (len([]rune(block.BlockStop)) - 1))]),
				}
				for k, v := range block.InjectValues {
					data[k] = v
				}
				dataSet = append(dataSet, data)
				activeBlock = nil
				p.resetBlocks(instructions, nil)
				nonMatchStart = i + 1
			} else if blockIsActive(block) && len([]rune(block.BlockStop)) > block.MatchIndex && []rune(block.BlockStop)[block.MatchIndex] == char {
				block.MatchIndex++
			} else if blockIsActive(block) {
				block.MatchIndex = 0
			}
			if activeBlock != nil {
				continue
			}

			// When an active block is not selected yet
			if len([]rune(block.BlockStart)) > block.MatchIndex && []rune(block.BlockStart)[block.MatchIndex] == char {
				block.Matching = true
				block.MatchIndex++
				if (block.MatchIndex == len([]rune(block.BlockStart))) && (candidateBlock == nil || block.MatchIndex > candidateBlock.MatchIndex) {
					candidateBlock = block
				}
			}
		}
		if candidateBlock != nil {
			activeBlock = candidateBlock
			candidateBlock = nil
			p.resetBlocks(instructions, activeBlock)
			nonMatchEnd := (i + 1) - len(activeBlock.BlockStart)
			if nonMatchEnd-nonMatchStart > 0 && defaultBlock != nil {
				data := map[string]any{
					"Content": string([]rune(markup)[nonMatchStart:nonMatchEnd]),
				}
				for k, v := range defaultBlock.InjectValues {
					data[k] = v
				}
				dataSet = append(dataSet, data)
			}
		}
	}
	// remainder
	if nonMatchStart < len([]rune(markup)) && defaultBlock != nil {
		endIndex := index
		if killSwitch {
			endIndex -= len([]rune(fullQuit))
		}
		if nonMatchStart < endIndex {
			data := map[string]any{
				"Content": string([]rune(markup)[nonMatchStart:endIndex]),
			}
			for k, v := range defaultBlock.InjectValues {
				data[k] = v
			}
			dataSet = append(dataSet, data)
		}
	}
	return dataSet, index
}

func (p *Parser) resetBlocks(instructions []*Block, exclude *Block) {
	for _, block := range instructions {
		block.MatchIndex = 0
		if !reflect.DeepEqual(block, exclude) {
			block.ContentStartIndex = 0
			block.Matching = false
		}
	}
}

func (p *Parser) parseKeys(text string) map[string]string {
	resp := map[string]string{}
	keySlice := strings.Split(text, "<")
	if len(keySlice) < 2 {
		return resp
	}
	// in case "<" is a char in the actual content
	keyString := strings.Join(keySlice[1:], "<")
	if []rune(keyString)[len(text)-1] == '>' {
		keyPairs := strings.Split(string([]rune(keyString)[:len(keyString)-1]), ",")
		for _, keyPair := range keyPairs {
			if keyVal := strings.Split(keyPair, ":"); len(keyVal) == 2 {
				resp[strings.TrimSpace(keyVal[0])] = strings.TrimSpace(keyVal[1])
			}
		}
	}
	return resp
}

var DefaultInstructions []*Block = []*Block{
	{
		BlockStart: "//p\n",
		BlockStop:  "\n//p",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "p",
		},
	},
	{
		BlockStart: "# ",
		BlockStop:  "\n",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "h1",
		},
	},
	{
		BlockStart: "## ",
		BlockStop:  "\n",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "h2",
		},
	},
	{
		BlockStart: "### ",
		BlockStop:  "\n",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "h3",
		},
	},
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
	{
		BlockStart: "<span>",
		BlockStop:  "</span>",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "span",
		},
	},
	{
		BlockStart: "//code-Go\n",
		BlockStop:  "\n//code",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type":     "code",
			"Language": "Go",
		},
	},
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
		},
	},
	{
		BlockStart: "//ul\n",
		BlockStop:  "\n//ul",
		SubBlocks: []*Block{
			{
				BlockStart: "- ",
				BlockStop:  "\n",
				SubBlocks:  nil,
				InjectValues: map[string]any{
					"Type": "li",
				},
			},
		},
		InjectValues: map[string]any{
			"Type": "ul",
		},
	},
	{
		BlockStart: "//ol\n",
		BlockStop:  "\n//ol",
		SubBlocks: []*Block{
			{
				BlockStart: "- ",
				BlockStop:  "\n",
				SubBlocks:  nil,
				InjectValues: map[string]any{
					"Type": "li",
				},
			},
		},
		InjectValues: map[string]any{
			"Type": "ol",
		},
	},
}
