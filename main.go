package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// todo
//
// recursion when subBlocks is not nil -> recall with remaining text with sub blocks as instruct and blockStop as fullQuit
//     should return index where we stop so caller of recursive func can skip ahead.
//
// block attributes and endblock tags -> inject values into article

type Block struct {
	KeyStart          string
	KeyStop           string
	BlockStop         string
	Matching          bool
	MatchIndex        int
	Subdir            map[string]*Block
	TypeName          string
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
// 
// \\ol
// - one
// - two
// - three
// \\ol
// `

// inlinep
// t := `
// //p
// Once I was **cool** I was so [cool] that I was a \cool\ boy!
// //p
// `


func main() {

	p := Block{
		KeyStart:  "//p\n",
		BlockStop: "\n//p",
		Subdir:    nil,
		TypeName:  "p",
	}
	h1 := Block{
		KeyStart:  "# ",
		BlockStop: "\n",
		Subdir:    nil,
		TypeName:  "h1",
	}
	h2 := Block{
		KeyStart:  "## ",
		BlockStop: "\n",
		Subdir:    nil,
		TypeName:  "h2",
	}
	h3 := Block{
		KeyStart:  "### ",
		BlockStop: "\n",
		Subdir:    nil,
		TypeName:  "h3",
	}

	instruct := []*Block{&h1, &h2, &h3, &p}
	fullQuit := "//e"
	Article, _ := doParse(instruct, md, fullQuit)
	r, _ := json.MarshalIndent(Article, "", "    ")
	fmt.Println(string(r))
}

func doParse(instructions []*Block, text string, fullQuit string) ([]map[string]any, int) {
	
	Article := []map[string]any{}
	fullQuitMatch := 0
	var activeBlock *Block
	var index int
	
	blockIsActive := func (block *Block) bool {
		return reflect.DeepEqual(activeBlock, block) 
	}

	var candidateBlock *Block = nil
	for i, char := range md {

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
			
			// when an active block is selected
			if activeBlock != nil && !blockIsActive(block){
				continue
			} else if activeBlock != nil && block.ContentStartIndex == 0 {
				block.ContentStartIndex = i
			}
			if blockIsActive(block) && len([]rune(block.BlockStop)) == (block.MatchIndex+1) && []rune(block.BlockStop)[block.MatchIndex] == char {
				// append to Article
				data := map[string]any{
					"Type":    block.TypeName,
					"Content": string([]rune(md[block.ContentStartIndex:(i - (len([]rune(block.BlockStop)) - 1))])),
				}
				// if there are subblocks, recur against remaining string with subblocks as instruct and blockStop as fullQuit
				// then add result as content and skip ahead to returned index
				Article = append(Article, data)
				activeBlock = nil
				resetBlocks(instructions, nil)
			} else if blockIsActive(block) && len([]rune(block.BlockStop)) > block.MatchIndex && []rune(block.BlockStop)[block.MatchIndex] == char {
				block.MatchIndex++
			} 
			if activeBlock != nil {
				continue
			}

			// When an active block is not selected yet
			if len([]rune(block.KeyStart)) > block.MatchIndex && []rune(block.KeyStart)[block.MatchIndex] == char {
				block.Matching = true
				block.MatchIndex++
				if (block.MatchIndex == len([]rune(block.KeyStart))) && (candidateBlock == nil || block.MatchIndex > candidateBlock.MatchIndex) {
					candidateBlock = block
				}
			}
		}
		if candidateBlock != nil {
			activeBlock = candidateBlock
			candidateBlock = nil
			resetBlocks(instructions, activeBlock)
		}
		index++
	}
	return Article, index
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
