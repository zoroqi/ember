package ember

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type Block struct {
	num   int
	Book  string
	Index string
	Text  string
	Err   error
}

func isBlockStart(s string) bool {
	if s == "" {
		return false
	}
	for _, ss := range s {
		if ss != '=' {
			return false
		}
	}
	return true
}

// 我也不知道为啥有会有这种情况
func removeSpecialCharset(s string) string {
	s = strings.TrimSpace(s)
	return strings.ReplaceAll(s, "\ufeff", "")
}

func ParseClippings(clippingsText io.Reader) (clips []Clipping, errBlock []Block) {
	blocks := ParseBlocks(clippingsText)
	for _, v := range blocks {
		c, err := BlockToClipping(v)
		if err != nil {
			v.Err = err
			errBlock = append(errBlock, v)
		} else {
			clips = append(clips, c)
		}

	}
	return
}

func ParseBlocks(clippingsText io.Reader) (blocks []Block) {
	lines := bufio.NewReader(clippingsText)
	var clip Block
	i := 0
	lineNum := 0
	for {
		lbs, _, err := lines.ReadLine()
		if err != nil {
			break
		}
		lineNum++
		line := removeSpecialCharset(string(lbs))
		if isBlockStart(line) {
			blocks = append(blocks, clip)
			clip = Block{num: lineNum}
			i = 0
			continue
		}
		switch {
		case i == 0:
			clip.Book = line
		case i == 1:
			clip.Index = line[2:]
		case i == 3:
			clip.Text = line
		case i > 3:
			clip.Text = clip.Text + "\n" + line
		}
		i++
	}
	if len(blocks) > 1 {
		if blocks[len(blocks)-1].num != lineNum {
			if clip.Index != "" {
				blocks = append(blocks, clip)
			}
		}
	}
	return
}

type ClippingType string

// 标记
const ClippingHighlight ClippingType = "highlight"

// 笔记
const ClippingNote ClippingType = "note"

// 书签
const ClippingBookmark ClippingType = "bookmark"

type ClippingLanguage interface {
	parseBlock(Block) (Clipping, error)
	String() string
}

type Clipping struct {
	Book          string
	Text          string
	Time          time.Time
	LocationStart int
	LocationEnd   int
	Type          ClippingType
}

func BlockToClipping(b Block) (Clipping, error) {
	cLanguage := language(b)
	c, err := cLanguage.parseBlock(b)
	return c, err
}

func language(s Block) ClippingLanguage {
	if strings.HasPrefix(s.Index, "您在") || strings.HasPrefix(s.Index, "在位置") {
		c := CLZhs("")
		return &c
	}

	if strings.HasPrefix(s.Index, "Your ") {
		c := CLEnglish("")
		return &c
	}

	c := CLEnglish("")
	return &c
}
