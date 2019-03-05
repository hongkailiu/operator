package lc

import "fmt"

func d(t *node) int {
	if t == nil {
		return 0
	}
	dl := d(t.l)
	dr := d(t.r)
	if dl < dr {
		return dr + 1
	}
	return dl + 1
}

func printT(prefix string, t *node, isLeft bool) {
	if t != nil {
		var s string
		if isLeft {
			s = fmt.Sprintf("%s|-- %d", prefix, t.v)
		} else {
			s = fmt.Sprintf("%s\\-- %d", prefix, t.v)
		}
		fmt.Println(s)

		if isLeft {
			s = fmt.Sprintf("%s|   ", prefix)
		} else {
			s = fmt.Sprintf("%s    ", prefix)
		}

		printT(s, t.l, true)
		printT(s, t.r, false)
	}
}
