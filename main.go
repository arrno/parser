package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

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

var md string = `### H
# H
## H
### H
## H
# H
## H
## H

//p
one I was a cool
man who was a cool man
cool
//p

//inlinep
Once I was **cool** I was so [[cool]] that I was a __cool__ boy!
//inlinep

//code-Go
fmt.Println("Hello")
//code

//e
## I will be cut out
`

// lists
// t := `
// \\ul
// - one
// - two
// - three
// \\ul

// \\ol
// - one
// - two
// - three
// \\ol
// `

func main() {

	p := Block{
		BlockStart: "//p\n",
		BlockStop:  "\n//p",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "p",
		},
	}
	h1 := Block{
		BlockStart: "# ",
		BlockStop:  "\n",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "h1",
		},
	}
	h2 := Block{
		BlockStart: "## ",
		BlockStop:  "\n",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "h2",
		},
	}
	h3 := Block{
		BlockStart: "### ",
		BlockStop:  "\n",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "h3",
		},
	}
	goCode := Block{
		BlockStart: "//code-Go\n",
		BlockStop:  "\n//code",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type":     "code",
			"Language": "Go",
		},
	}
	bold := Block{
		BlockStart: "**",
		BlockStop:  "**",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "b",
		},
	}
	italics := Block{
		BlockStart: "__",
		BlockStop:  "__",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "i",
		},
	}
	inlineCode := Block{
		BlockStart: "[[",
		BlockStop:  "]]",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "code",
		},
	}
	span := Block{
		BlockStart: "<span>",
		BlockStop:  "</span>",
		SubBlocks:  nil,
		InjectValues: map[string]any{
			"Type": "span",
		},
	}
	inlinep := Block{
		BlockStart: "//inlinep\n",
		BlockStop:  "\n//inlinep",
		SubBlocks: []*Block{
			&bold,
			&italics,
			&inlineCode,
		},
		SubDefaultBlock: &span,
		InjectValues: map[string]any{
			"Type": "inlinep",
		},
	}

	instructions := []*Block{&h1, &h2, &h3, &p, &goCode, &inlinep}
	fullQuit := "//e"
	DataSet, _ := doParse(instructions, md, fullQuit, nil)
	r, _ := json.MarshalIndent(DataSet, "", "    ")
	fmt.Println(string(r))
}

func doParse(instructions []*Block, text string, fullQuit string, defaultBlock *Block) ([]map[string]any, int) {

	dataSet := []map[string]any{}
	fullQuitMatch := 0
	var activeBlock *Block
	var index int

	blockIsActive := func(block *Block) bool {
		return reflect.DeepEqual(activeBlock, block)
	}

	nonMatchStart := 0
	var candidateBlock *Block = nil
	skipTo := 0

	for i, char := range text {
		index++

		// handle for recursion
		if []rune(fullQuit)[fullQuitMatch] == char {
			fullQuitMatch++
		} else {
			fullQuitMatch = 0
		}
		if len([]rune(fullQuit)) == fullQuitMatch {
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
				content, skipBy = doParse(block.SubBlocks, string([]rune(text)[i:len([]rune(text))]), block.BlockStop, block.SubDefaultBlock)
				skipTo = i + skipBy
				data := map[string]any{
					"Content": content,
				}
				for k, v := range block.InjectValues {
					data[k] = v
				}
				dataSet = append(dataSet, data)
				activeBlock = nil
				resetBlocks(instructions, nil)
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
					"Content": string([]rune(text)[block.ContentStartIndex:(i - (len([]rune(block.BlockStop)) - 1))]),
				}
				for k, v := range block.InjectValues {
					data[k] = v
				}
				dataSet = append(dataSet, data)
				activeBlock = nil
				resetBlocks(instructions, nil)
				nonMatchStart = i + 1
			} else if blockIsActive(block) && len([]rune(block.BlockStop)) > block.MatchIndex && []rune(block.BlockStop)[block.MatchIndex] == char {
				block.MatchIndex++
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
			resetBlocks(instructions, activeBlock)
			nonMatchEnd := (i + 1) - len(activeBlock.BlockStart)
			if nonMatchEnd-nonMatchStart > 0 && defaultBlock != nil {
				data := map[string]any{
					"Content": string([]rune(text)[nonMatchStart:nonMatchEnd]),
				}
				for k, v := range defaultBlock.InjectValues {
					data[k] = v
				}
				dataSet = append(dataSet, data)
			}
		}
	}
	// remainder
	if nonMatchStart < len([]rune(text)) && defaultBlock != nil {
		data := map[string]any{
			// TODO this currently includes the quitTag...
			"Content": string([]rune(text)[nonMatchStart:index]),
		}
		for k, v := range defaultBlock.InjectValues {
			data[k] = v
		}
		dataSet = append(dataSet, data)
	}
	return dataSet, index
}

// An idea
func parseStopKeys(blockStop string) map[string]string {
	resp := map[string]string{}
	keySlice := strings.Split(blockStop, "<")
	if len(keySlice) < 2 {
		return resp
	}
	// in case "<" is a char in the actual content
	keyString := strings.Join(keySlice[1:], "<")
	if []rune(keyString)[len(blockStop)-1] == '>' {
		keyPairs := strings.Split(string([]rune(keyString)[:len(keyString)-1]), ",")
		for _, keyPair := range keyPairs {
			if keyVal := strings.Split(keyPair, ":"); len(keyVal) == 2 {
				resp[strings.TrimSpace(keyVal[0])] = strings.TrimSpace(keyVal[1])
			}
		}
	}
	return resp
}

func resetBlocks(blocks []*Block, exclude *Block) {
	for _, block := range blocks {
		block.MatchIndex = 0
		if !reflect.DeepEqual(block, exclude) {
			block.ContentStartIndex = 0
			block.Matching = false
		}
	}
}
