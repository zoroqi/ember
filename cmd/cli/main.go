package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zoroqi/ember"
	"os"
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
		fmt.Println(err)
		return
	}
	defer reader.Close()
	clips, errBlock := ember.ParseClippings(reader)
	bs, err := json.MarshalIndent(clips, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bs))

	if len(errBlock) != 0 {
		bs, err := json.MarshalIndent(errBlock, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", string(bs))
		}
		fmt.Fprintf(os.Stderr, "%s", string(bs))
	}
}
