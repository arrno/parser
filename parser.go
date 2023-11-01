package parser

import (
	"reflect"
	"strings"
)

// MarkupParser is something that is able to parse markup text into structured data.
type MarkupParser interface {
	DoParse(markup string) []map[string]any
	ParseMapKeys(text string) (result map[string]string, parsed int)
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
				keyVals, parsed := p.ParseMapKeys(markup[min(skipTo, len([]rune(markup))-1):])
				if parsed > 0 {
					skipTo += parsed
					for k, v := range keyVals {
						data[k] = v
					}
				}
				dataSet = append(dataSet, data)
				activeBlock = nil
				p.resetBlocks(instructions, nil)
				nonMatchStart = skipTo + parsed
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
				// 10 3
				keyVals, parsed := p.ParseMapKeys(markup[min(i+1, len([]rune(markup))-1):])
				if parsed > 0 {
					skipTo = i + parsed + 1
					for k, v := range keyVals {
						data[k] = v
					}
				}
				dataSet = append(dataSet, data)
				activeBlock = nil
				p.resetBlocks(instructions, nil)
				nonMatchStart = i + 1 + parsed
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

type inheritedContent struct {
	content      []map[string]any
	openTagStart int
	closeTagEnd  int
}
type stackBlock struct {
	openTagStart     int
	closeTagEnd      int
	matchIndex       int
	inheritedContent inheritedContent
	block            *Block
}

// stack method for infinitely deep markup
func (p *Parser) handleParseStack(instructions []*Block, markup string) []map[string]any {

	dataSet := []map[string]any{}

	activeBlockStack := []*stackBlock{}
	var candidateBlock *Block = nil

	blockIsActive := func(block *Block) bool {
		return len(activeBlockStack) > 0 && reflect.DeepEqual(activeBlockStack[len(activeBlockStack)-1].block, block)
	}

	skipBy := 0

	for i, char := range markup {

		for _, block := range instructions {

			if skipBy > 0 {
				skipBy--
				continue
			}

			// evaluate for pushing onto stack
			if len([]rune(block.BlockStart)) > block.MatchIndex && []rune(block.BlockStart)[block.MatchIndex] == char {
				block.MatchIndex++
				if (block.MatchIndex == len([]rune(block.BlockStart))) && (candidateBlock == nil || block.MatchIndex > candidateBlock.MatchIndex) {
					candidateBlock = block
				}
			}

			// evaluate for popping off of stack
			if blockIsActive(block) && len([]rune(block.BlockStop)) > block.MatchIndex && []rune(block.BlockStop)[block.MatchIndex] == char {
				block.MatchIndex++
				if block.MatchIndex == len([]rune(block.BlockStop)) {
					// pull text between active block tags
					// check for inherited content... put matching text between inherited content into defaul tags and merge into inherited content
					// parseMapKeys, set skip note closeTagEnd
					// if this block is at position 1, append merged content into dataSet
					//     else, inject as inherited content on next block in stack
					// pop
					activeBlockStack = activeBlockStack[:len(activeBlockStack)-1]
				}
			}
		}
		// execute push onto stack
		if candidateBlock != nil {
			sb := stackBlock{
				openTagStart: i - len([]rune(candidateBlock.BlockStart)),
				closeTagEnd:  -1,
				matchIndex:   0,
				block:        candidateBlock,
			}
			activeBlockStack = append(activeBlockStack, &sb)
			candidateBlock = nil
			p.resetBlocks(instructions, nil)
		}
	}

	return dataSet
}

// ParseMapKeys takes in a string and attempts to parse a key value pair
// pattern at the beginning. If successful, the number of runes parsed is also
// returned.
//
// expectation is ::[key: val, key2: val2, keyn: valn] ...
func (p *Parser) ParseMapKeys(text string) (result map[string]string, parsed int) {
	var matchString string
	result = map[string]string{}
	runeSlice := []rune(text)
	if len(runeSlice) < 7 || string(runeSlice[:3]) != "::[" {
		return
	} else if temp := strings.Split(text[3:], "]"); len(temp) < 2 {
		return
	} else {
		matchString = temp[0]
	}
	keyPairs := strings.Split(string(matchString), ",")
	for _, keyPair := range keyPairs {
		if keyVal := strings.Split(keyPair, ":"); len(keyVal) == 2 {
			result[strings.TrimSpace(keyVal[0])] = strings.TrimSpace(keyVal[1])
		}
	}
	parsed = len([]rune(matchString)) + 4
	return
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
