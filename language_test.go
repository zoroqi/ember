package ember

import "testing"

func TestDefaultLocation(t *testing.T) {
	d := map[string][]int{}
	d["您在第 57 页的笔记"] = []int{0, 57}
	d["您在位置 #1680 的书签"] = []int{0, 1680}
	d["您在第 126-126 页的标注"] = []int{126, 126}
	d["您在第 57 页的笔记"] = []int{0, 57}
	d["您在位置 #1680 的书签"] = []int{0, 1680}
	d["您在位置 #1103-1103的标注"] = []int{1103, 1103}
	d["您在位置 #1385 的笔记"] = []int{0, 1385}
	d["Your Highlight on Location 130-131"] = []int{130, 131}
	d["Your Note on Location 948"] = []int{0, 948}
	for k, v := range d {
		start, end, err := defaultLocation(k)
		if err != nil {
			t.Error("error", k, v, start, end, err)
		}
		if start != v[0] || end != v[1] {
			t.Error("parse error", k, v, start, end, err)
		}
	}

}
