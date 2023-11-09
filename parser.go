package parser

import (
	"reflect"
	"strings"
)

const (
	defaultTrimChars string = " \n	"
)

type trimOpt uint

const (
	trimLeft trimOpt = iota
	trimRight
	trimAll
	trimNewLine
)

// MarkupParser is something that is able to parse markup text into structured data.
type MarkupParser interface {
	DoParse(markup string) []map[string]any
	ParseMapKeys(text string) (result map[string]string, parsed int)
	SetTrimChars(text string)
}

// Instructions are meant to be written by the user.
type Instruction struct {
	OpenTag    string
	CloseTag   string
	Attributes map[string]any
	Trim       bool
}

// ParseInstructions are a set of Instructions used to instantiate a parser.
type ParseInstructions []Instruction

// inheritedContent is for block level tracking of what sub text has already been parsed.
type inheritedContent struct {
	content    []map[string]any
	startStops [][]int
}

// stackBlock is what is pushed/popped from the stack.
type stackBlock struct {
	openTagStart     int
	matchIndex       int
	inheritedContent *inheritedContent
	block            *block
}

// block is a single markup parsing instruction. A block can have nested sub blocks.
type block struct {
	BlockStart     string
	BlockStop      string
	MatchIndex     int
	MatchStopIndex int
	InjectValues   map[string]any
	Trim           bool
}

// Parser implements MarkupParser interface
type Parser struct {
	instructions []*block
	trimChars    string
}

// NewParser creates a new Parser and returns a pointer to it.
func NewParser(instructions ParseInstructions) *Parser {
	if instructions == nil {
		instructions = DefaultInstructions
	}
	blocks := make([]*block, len(instructions))
	for i, instruction := range instructions {
		blocks[i] = &block{
			BlockStart:     instruction.OpenTag,
			BlockStop:      instruction.CloseTag,
			MatchIndex:     0,
			MatchStopIndex: 0,
			InjectValues:   instruction.Attributes,
			Trim:           instruction.Trim,
		}
	}
	p := Parser{
		instructions: blocks,
		trimChars:    defaultTrimChars,
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

// stack method for infinitely deep markup
func (p *Parser) handleParseStack(markup string) []map[string]any {

	dataSet := []map[string]any{}
	markupRune := []rune(markup)

	activeBlockStack := []*stackBlock{}
	var candidateBlock *block = nil

	blockIsActive := func(block *block) bool {
		return len(activeBlockStack) > 0 && reflect.DeepEqual(activeBlockStack[len(activeBlockStack)-1].block, block)
	}

	skipBy := 0

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
						for idx, startStop := range stackBlock.inheritedContent.startStops {
							// if unmatch start is less than startStop[0] -> append unmatchStart:0, append inherited, unmatchStart = startStop[1]
							text := string(markupRune[unmatchStart:startStop[0]])
							if unmatchStart < startStop[0] {
								if block.Trim && idx == 0 {
									text = p.handleTrim(text, trimLeft)
								} else if block.Trim {
									text = p.handleTrim(text, trimNewLine)
								}
								subData := map[string]any{
									"Content": text,
								}
								if len(text) > 0 {
									contentSlice = append(contentSlice, subData)
								}
							}
							unmatchStart = startStop[1] + 1
							contentSlice = append(contentSlice, stackBlock.inheritedContent.content[idx])
						}
						if unmatchStart < i-(len(block.BlockStop)) {
							text := string(markupRune[unmatchStart : i-(len(blockStartRune))])
							var subData map[string]any
							if block.Trim {
								text = p.handleTrim(text, trimRight)
							}
							subData = map[string]any{
								"Content": text,
							}
							if len(text) > 0 {
								contentSlice = append(contentSlice, subData)
							}
						}
						data["Content"] = contentSlice
					} else {
						text := string(matchTextRune)
						if block.Trim {
							text = p.handleTrim(text, trimAll)
						}
						data["Content"] = text
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
			} else if blockIsActive(block) {
				block.MatchStopIndex = 0
			}

			// evaluate for pushing onto stack
			if len(blockStartRune) > block.MatchIndex && blockStartRune[block.MatchIndex] == char {
				block.MatchIndex++
				if (block.MatchIndex == len(blockStartRune)) && (candidateBlock == nil || block.MatchIndex > candidateBlock.MatchIndex) {
					candidateBlock = block
				}
			} else {
				block.MatchIndex = 0
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

func (p *Parser) handleTrim(text string, opt trimOpt) string {
	lines := strings.Split(text, "\n")
	nonEmptyLines := []string{}
	// for _, line := range lines {
	// 	if len(strings.Trim(line, p.trimChars)) > 0 {
	// 		nonEmptyLines = append(nonEmptyLines, line)
	// 	}
	// }
	for i, line := range lines {
		if opt == trimAll {
			line = strings.Trim(line, p.trimChars)
		} else {
			// only trim right excludes left most trim
			if opt == trimLeft || i > 0 {
				line = strings.TrimLeft(line, p.trimChars)
			}
			// only trim right excludes right most trim
			if opt == trimRight || i < len(nonEmptyLines)-1 {
				line = strings.TrimRight(line, p.trimChars)
			}
		}
		if len(line) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return strings.Join(nonEmptyLines, " ")
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

// SetTrimChars updates the set of characters that are trimmed from applicable blocks.
//
// The input variable represents a set of runes rather than a literal full match.
func (p *Parser) SetTrimChars(text string) {
	p.trimChars = text
}
