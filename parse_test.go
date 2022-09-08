package ember

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParseClippings(t *testing.T) {
	data, _ := os.Open("test_data.txt")
	defer data.Close()
	clips, err := ParseClippings(data)
	if len(err) != 0 {
		t.Fatal(err)
	}

	mustTimeParser := func(s string) time.Time {
		t, e := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		if e != nil {
			panic(e)
		}
		return t
	}

	book := "test-book (ember)"
	var testData []Clipping
	testData = append(testData, Clipping{
		Book:          book,
		LocationStart: 1, LocationEnd: 2, Text: "test1", Type: ClippingHighlight,
		Time: mustTimeParser("2022-05-25 10:05:05"),
	})
	testData = append(testData, Clipping{
		Book:          book,
		LocationStart: 0, LocationEnd: 3, Text: "test2", Type: ClippingNote,
		Time: mustTimeParser("2022-06-02 09:15:59"),
	})
	testData = append(testData, Clipping{
		Book:          book,
		LocationStart: 4, LocationEnd: 5, Text: "test3", Type: ClippingHighlight,
		Time: mustTimeParser("2022-06-06 09:52:14"),
	})
	testData = append(testData, Clipping{
		Book:          book,
		LocationStart: 0, LocationEnd: 6, Text: "test4", Type: ClippingNote,
		Time: mustTimeParser("2022-06-07 09:25:11"),
	})
	testData = append(testData, Clipping{
		Book:          book,
		LocationStart: 0, LocationEnd: 6, Text: "test5\nnewline", Type: ClippingNote,
		Time: mustTimeParser("2022-06-07 09:25:11"),
	})

	for i, c := range clips {
		if !reflect.DeepEqual(c, testData[i]) {
			t.Fatalf("parse error, %v, but %v", testData[i], c)
		}
	}
}
