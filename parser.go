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
	BlockStart     string
	BlockStop      string
	MatchIndex     int
	MatchStopIndex int
	InjectValues   map[string]any
}

// Parser implements MarkupParser interface
type Parser struct {
	instructions []*Block
}

// NewParser creates a new Parser and returns a pointer to it.
func NewParser(instructions ParseInstructions) *Parser {
	if instructions == nil {
		instructions = DefaultInstructions
	}
	blocks := make([]*Block, len(instructions))
	for i, instruction := range instructions {
		blocks[i] = &Block{
			BlockStart:     instruction.OpenTag,
			BlockStop:      instruction.CloseTag,
			MatchIndex:     0,
			MatchStopIndex: 0,
			InjectValues:   instruction.Attributes,
		}
	}
	p := Parser{
		instructions: blocks,
	}
	return &p
}

// DoParse will parse the input markup text into structured data.
func (p *Parser) DoParse(markup string) []map[string]any {
	return p.handleParseStack(markup)
}

func (p *Parser) resetBlocks() {
	for _, block := range p.instructions {
		block.MatchIndex = 0
		block.MatchStopIndex = 0
	}
}

type Instruction struct {
	OpenTag    string
	CloseTag   string
	Attributes map[string]any
}

type ParseInstructions []Instruction
type inheritedContent struct {
	content    []map[string]any
	startStops [][]int
}
type stackBlock struct {
	openTagStart     int
	matchIndex       int
	inheritedContent *inheritedContent
	block            *Block
}

// stack method for infinitely deep markup
func (p *Parser) handleParseStack(markup string) []map[string]any {

	dataSet := []map[string]any{}
	markupRune := []rune(markup)

	activeBlockStack := []*stackBlock{}
	var candidateBlock *Block = nil

	blockIsActive := func(block *Block) bool {
		return len(activeBlockStack) > 0 && reflect.DeepEqual(activeBlockStack[len(activeBlockStack)-1].block, block)
	}

	skipBy := 0

	// for byteIndex, char := range markup {
	for i := 0; i < len(markupRune); i++ {

		char := markupRune[i]
		for _, block := range p.instructions {

			if skipBy > 0 {
				skipBy--
				continue
			}

			blockStopRune := []rune(block.BlockStop)
			blockStartRune := []rune(block.BlockStart)

			// evaluate for popping off of stack
			if blockIsActive(block) && len(blockStopRune) > block.MatchStopIndex && blockStopRune[block.MatchStopIndex] == char {
				block.MatchStopIndex++
				if block.MatchStopIndex == len(blockStopRune) {
					stackBlock := activeBlockStack[len(activeBlockStack)-1]
					// pull text between active block tags
					matchTextRune := markupRune[(stackBlock.openTagStart + len(blockStartRune)) : i-(len(blockStopRune)-1)]
					data := map[string]any{}
					if stackBlock.inheritedContent != nil {
						contentSlice := []map[string]any{}
						// check for inherited content... put matching text between inherited content into defaul tags and merge into inherited content
						unmatchStart := (stackBlock.openTagStart + len(blockStartRune))
						for i, startStop := range stackBlock.inheritedContent.startStops {
							// if unmatch start is less than startStop[0] -> append unmatchStart:0, append inherited, unmatchStart = startStop[1]
							if unmatchStart < startStop[0] {
								subData := map[string]any{
									"Content": string(markupRune[unmatchStart:startStop[0]]),
								}
								contentSlice = append(contentSlice, subData)
							}
							unmatchStart = startStop[1] + 1
							contentSlice = append(contentSlice, stackBlock.inheritedContent.content[i])
						}
						if unmatchStart < i-(len(block.BlockStop)-1) {
							subData := map[string]any{
								"Content": string(markupRune[unmatchStart : i-(len(blockStartRune)-1)]),
							}
							contentSlice = append(contentSlice, subData)
						}
						data["Content"] = contentSlice
					} else {
						data["Content"] = string(matchTextRune)
					}
					for k, v := range block.InjectValues {
						data[k] = v
					}
					// parseMapKeys, set skip note closeTagEnd
					keyVals, parsed := p.ParseMapKeys(string(markupRune[min(i+1, len(markupRune)-1):]))
					if parsed > 0 {
						skipBy += parsed
						for k, v := range keyVals {
							data[k] = v
						}
					}
					// if this block is at position 1, append merged content into dataSet
					if len(activeBlockStack) == 1 {
						dataSet = append(dataSet, data)
					} else {
						// else, inject as inherited content on next block in stack
						ic := activeBlockStack[len(activeBlockStack)-2].inheritedContent
						if ic == nil {
							ic = &inheritedContent{
								content:    []map[string]any{},
								startStops: [][]int{},
							}
						}
						ic.content = append(ic.content, data)
						ic.startStops = append(ic.startStops, []int{stackBlock.openTagStart, i + parsed})
						activeBlockStack[len(activeBlockStack)-2].inheritedContent = ic
					}
					// pop
					activeBlockStack = activeBlockStack[:len(activeBlockStack)-1]
					p.resetBlocks()
				}
			}

			// evaluate for pushing onto stack
			if len(blockStartRune) > block.MatchIndex && blockStartRune[block.MatchIndex] == char {
				block.MatchIndex++
				if (block.MatchIndex == len(blockStartRune)) && (candidateBlock == nil || block.MatchIndex > candidateBlock.MatchIndex) {
					candidateBlock = block
				}
			}
		}
		// execute push onto stack
		if candidateBlock != nil {
			sb := stackBlock{
				openTagStart: i - len([]rune(candidateBlock.BlockStart)) + 1,
				matchIndex:   0,
				block:        candidateBlock,
			}
			activeBlockStack = append(activeBlockStack, &sb)
			candidateBlock = nil
			p.resetBlocks()
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

var DefaultInstructions ParseInstructions = ParseInstructions{
	{
		OpenTag:  "<p>",
		CloseTag: "</p>",
		Attributes: map[string]any{
			"Type": "p",
		},
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
	},
	{
		OpenTag:  "<h1>",
		CloseTag: "</h1>",
		Attributes: map[string]any{
			"Type": "h1",
		},
	},
	{
		OpenTag:  "<h2>",
		CloseTag: "</h2>",
		Attributes: map[string]any{
			"Type": "h2",
		},
	},
	{
		OpenTag:  "<h3>",
		CloseTag: "</h3>",
		Attributes: map[string]any{
			"Type": "h3",
		},
	},
	{
		OpenTag:  "<div>",
		CloseTag: "</div>",
		Attributes: map[string]any{
			"Type": "d",
		},
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
	},
	{
		OpenTag:  "<li>",
		CloseTag: "</li>",
		Attributes: map[string]any{
			"Type": "li",
		},
	},
}
