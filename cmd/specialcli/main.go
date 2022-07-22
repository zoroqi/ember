package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/zoroqi/ember"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	filePath := flag.String("f", "", "kindle 'My Clippings.txt' path")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("file path is empty")
		return
	}
	reader, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("open file err %s, %v", *filePath, err)
		return
	}
	defer reader.Close()
	books, errBlocks := ember.ParseClippings(reader)

	if len(errBlocks) != 0 {
		for _, v := range errBlocks {
			fmt.Printf("%v\n", v)
		}
		return
	}

	bookCache := buildCache(books)

	help := `
Can only appear once
b list all books
f find book name
s list specify book by id
`
	for true {
		fmt.Println(help)
		input := strings.TrimSpace(readConsoleLine())
		if "" == input {
			continue
		}

		param := strings.SplitN(input, " ", 2)

		switch param[0] {
		case "b":
			bookCache.booksList(os.Stdout)
		case "f":
			if len(param) <= 1 {
				fmt.Println("input keyword")
				continue
			}
			bookCache.findBook(param[1], os.Stdout)
		case "s":
			{
				if len(param) <= 1 {
					fmt.Println("input book id")
					continue
				}
				id, err := strconv.Atoi(param[1])
				if err != nil {
					fmt.Println("id is not number")
					continue
				}
				bookCache.printClips(id, os.Stdout)
			}
		default:
			fmt.Println("error param")

		}

	}
}

func readConsoleLine() string {
	reader := bufio.NewReader(os.Stdin)
	data, _, e := reader.ReadLine()
	if e != nil {
		return ""
	}
	regexStr := string(data)
	return regexStr
}

type book struct {
	Id   int
	Book string
	Time time.Time
	Date string
	Text string
	Note string
}

func buildCache(clips []ember.Clipping) cache {

	books := make([]book, 0, len(clips))
	notes := make([]ember.Clipping, 0, len(clips)/50)
	offsetMapping := make(map[string]map[int]int)

	for i, c := range clips {
		switch c.Type {
		case ember.ClippingHighlight:
			books = append(books, book{Id: i, Book: c.Book, Time: c.Time, Text: c.Text, Note: ""})
			if _, exist := offsetMapping[c.Book]; !exist {
				offsetMapping[c.Book] = make(map[int]int)
			}
			offsetMapping[c.Book][c.LocationEnd] = i
		case ember.ClippingNote:
			notes = append(notes, c)
		}
	}

	for _, n := range notes {
		if o, exist := offsetMapping[n.Book]; exist {
			id := o[n.LocationEnd]
			if id > 0 {
				index := sort.Search(len(books), func(i int) bool {
					return books[i].Id >= id
				})
				if index >= 0 && index < len(books) {
					if books[index].Id == id {
						books[index].Note = n.Text
					}
				}
			}
		}
	}
	r := cache{clips: map[string][]book{}, bookName: []string{}}
	for _, v := range books {
		if _, ok := r.clips[v.Book]; !ok {
			r.bookName = append(r.bookName, v.Book)
		}
		r.clips[v.Book] = append(r.clips[v.Book], v)
	}
	return r
}

type cache struct {
	clips    map[string][]book
	bookName []string
}

func (c *cache) booksList(output io.Writer) {
	for i, v := range c.bookName {
		output.Write([]byte(fmt.Sprintf("%d. %s\n", i, v)))
	}
}

func (c *cache) findBook(find string, output io.Writer) {
	for i, v := range c.bookName {
		if strings.Contains(v, find) {
			output.Write([]byte(fmt.Sprintf("%d. %s\n", i, v)))
		}
	}
}

func (c *cache) printClips(id int, output io.Writer) {
	if id >= len(c.bookName) {
		return
	}
	clips := c.clips[c.bookName[id]]
	sa := make([]string, 0, len(clips))
	startReadIndex := 0
	saIndex := 0
	currentTime := clips[0].Time
	sa = append(sa, "")
	const interval = 86400 * 15
	for i, v := range clips {
		if (v.Time.Unix() - currentTime.Unix()) >= interval {
			sa[saIndex] = dateFormat(clips[startReadIndex], clips[i-1])
			startReadIndex = i
			sa = append(sa, "")
			saIndex = len(sa) - 1
		}
		sa = append(sa, format(i+1, v))
		currentTime = v.Time
	}
	sa[saIndex] = dateFormat(clips[startReadIndex], clips[len(clips)-1])

	output.Write([]byte(clips[0].Book))
	output.Write([]byte{'\n'})
	for _, v := range sa {
		output.Write([]byte(v))
	}
}

const time_layout = "2006-01-02"

func dateFormat(e1, e2 book) string {
	return fmt.Sprintf("* %s ~ %s\n", e1.Time.Format(time_layout), e2.Time.Format(time_layout))
}

func format(i int, e book) string {
	str := fmt.Sprintf("%d. %s\n", i, e.Text)
	if e.Note != "" {
		str += fmt.Sprintf("    * %s\n", e.Note)
	}
	return str
}
