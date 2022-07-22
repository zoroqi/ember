package ember

import (
	"regexp"
	"strings"
	"time"
)

func zhsDate(s string) (time.Time, error) {
	s = strings.ReplaceAll(s, "下午", "PM")
	s = strings.ReplaceAll(s, "上午", "AM")
	weekly := regexp.MustCompile("星期.")
	s = weekly.ReplaceAllString(s, "")
	timeFormat := "2006年1月2日 PM3:04:05"
	return time.ParseInLocation(timeFormat, s, time.Local)
}

func englishDate(s string) (time.Time, error) {
	timeFormat := "Monday, January 2, 2006 3:04:05 PM"
	return time.ParseInLocation(timeFormat, s, time.Local)
}
