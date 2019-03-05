package lc

import (
	"fmt"
	"strings"
)

func lengthOfLongestSubstring(s string) int {
	r := &result{input: s}
	r.init()
	r.split()
	return len(r.maxSub())
}

type result struct {
	input string
	m     map[rune]int
	subs  []string
	//r     string
}

func (r *result) init() {
	r.m = make(map[rune]int)
	for _, char := range r.input {
		r.m[char]++
	}
	r.subs = []string{r.input}
}

func (r *result) split() {
	if len(r.subs) == 0 {
		return
	}
	k, v := r.maxKV()
	r.m[k] = 0
	if v < 2 {
		//r.r = fmt.Sprintf("%s%s", r.r, r.maxSub())
		return
	}
	//r.r = fmt.Sprintf("%s%s", r.r, string(k))

	var newSubs []string
	for _, sub := range r.subs {
		newSubs = append(newSubs, getSubs(k, sub)...)
	}
	r.subs = newSubs
	//fmt.Println(fmt.Sprintf("===k,v: (=%s=,%d)", string(k), v))
	//fmt.Println(fmt.Sprintf("===subs: %v", r.subs))
	r.split()
}

func getSubs(r rune, s string) []string {
	words := strings.Split(s, string(r))
	if len(words) < 2 {
		return words
	}
	var result []string
	var pre string
	for i, w := range words {
		if i > 0 {
			result = append(result, fmt.Sprintf("%s%s%s", pre, string(r), w))
		}
		pre = w
	}
	return result
}

func (r *result) maxKV() (rune, int) {
	var rk rune
	var rv int
	for k, v := range r.m {
		if v > 1 && v > rv {
			rk = k
			rv = v
		}
	}
	return rk, rv
}

func (r *result) maxSub() string {
	var result string
	for _, sub := range r.subs {
		if len(sub) > len(result) {
			result = sub
		}
	}
	return result
}
