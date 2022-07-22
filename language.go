package ember

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var IndexError = errors.New("index error")
var NotFoundLocation = errors.New("not found location")

type CLZhs string

func (c CLZhs) parseBlock(b Block) (Clipping, error) {
	ss := strings.SplitN(b.Index, "|", 2)
	if len(ss) != 2 {
		return Clipping{}, IndexError
	}
	location := strings.TrimSpace(ss[0])
	var ct ClippingType
	switch {
	case strings.Contains(location, "笔记"):
		ct = ClippingNote
	case strings.Contains(location, "书签") || strings.Contains(location, "文章剪切"):
		ct = ClippingBookmark
	case strings.Contains(location, "标注"):
		ct = ClippingHighlight
	}

	// 您在第 390 页（位置 #5763-5765）的标注
	// 您在第 91 页（位置 #1347）的笔记
	// 这两种需要特殊处理, 产生的原因是 kindle 的 kfx 格式文件的标注
	if strings.Contains(location, "页（") {
		location = strings.Split(location, "页（")[1]
	}

	start, end, err := defaultLocation(location)
	if err != nil {
		return Clipping{}, err
	}

	dateStr := strings.TrimSpace(strings.ReplaceAll(ss[1], "添加于 ", ""))
	date, err := zhsDate(dateStr)
	if err != nil {
		return Clipping{}, err
	}
	return Clipping{
		Book:          b.Book,
		Text:          b.Text,
		LocationEnd:   end,
		LocationStart: start,
		Type:          ct,
		Time:          date,
	}, nil
}

func (c CLZhs) String() string {
	return "简体中文"
}

// - 您在第 57 页的笔记
// - 您在位置 #1680 的书签
// - 您在第 126-126 页的标注
// - 您在第 57 页的笔记
// - 您在位置 #1680 的书签
// - 您在位置 #1103-1103的标注
// - 您在位置 #1385 的笔记
// - Your Highlight on Location 941-948
// - Your Note on Location 948
var defaultLocationP = regexp.MustCompile("#?(\\d+)(-(\\d+))?")

func defaultLocation(s string) (start, end int, err error) {
	locations := defaultLocationP.FindStringSubmatch(s)
	var startS string
	var endS string
	startS = locations[1]
	endS = locations[3]
	if endS == "" {
		endS = startS
		startS = ""
	}
	if endS == "" && startS == "" {
		return 0, 0, NotFoundLocation
	}

	if startS != "" {
		start, err = strconv.Atoi(startS)
		if err != nil {
			err = fmt.Errorf("start location, %w", err)
			return
		}
	}
	if endS != "" {
		end, err = strconv.Atoi(endS)
		if err != nil {
			err = fmt.Errorf("end location, %w", err)
			return
		}
	}
	return
}

type CLEnglish string

// - Your Highlight on Location 941-948 | Added on Tuesday, June 7, 2022 10:21:11 AM
// - Your Note on Location 948 | Added on Tuesday, June 7, 2022 10:21:37 AM
func (c CLEnglish) parseBlock(b Block) (Clipping, error) {
	ss := strings.SplitN(b.Index, "|", 2)
	if len(ss) != 2 {
		return Clipping{}, IndexError
	}
	location := strings.TrimSpace(ss[0])
	var ct ClippingType
	switch {
	case strings.Contains(ss[0], "Note"):
		ct = ClippingNote
	case strings.Contains(ss[0], "Bookmark"):
		ct = ClippingBookmark
	case strings.Contains(ss[0], "Highlight"):
		ct = ClippingHighlight
	}
	start, end, err := defaultLocation(location)
	if err != nil {
		return Clipping{}, err
	}
	dateStr := strings.TrimSpace(strings.ReplaceAll(ss[1], "Added on ", ""))
	date, err := englishDate(dateStr)
	if err != nil {
		return Clipping{}, err
	}
	return Clipping{
		Book:          b.Book,
		Text:          b.Text,
		LocationEnd:   end,
		LocationStart: start,
		Type:          ct,
		Time:          date,
	}, nil
}

func (c CLEnglish) String() string {
	return "English"
}
