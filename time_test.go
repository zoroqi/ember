package ember

import (
	"testing"
	"time"
)

func TestZhsDate(t *testing.T) {
	zh := "2015年12月16日星期三 下午7:55:38"
	date, err := zhsDate(zh)
	if err != nil {
		t.Fatal(err)
	}
	d := time.Date(2015, 12, 16, 19, 55, 38, 0, time.Local)
	if date != d {
		t.Fatal("parse error")
	}
	zh = "2015年12月20日星期六 下午7:55:38"
	date, err = zhsDate(zh)
	if err != nil {
		t.Fatal(err)
	}
	d = time.Date(2015, 12, 20, 19, 55, 38, 0, time.Local)
	if date != d {
		t.Fatal("parse error")
	}
	zh = "2015年12月20日星期六 下午37:55:38"
	date, err = zhsDate(zh)
	if err != nil {
	}
}

func TestEnglishDate(t *testing.T) {
	zh := "Wednesday, June 8, 2022 10:48:16 AM"
	date, err := englishDate(zh)
	if err != nil {
		t.Fatal(err)
	}
	d := time.Date(2022, 6, 8, 10, 48, 16, 0, time.Local)
	if date != d {
		t.Fatal("parse error")
	}
	zh = "Wednesday, June 8, 2022 37:48:16 AM"
	date, err = zhsDate(zh)
	if err != nil {
	}
}
