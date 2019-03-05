package lc

import "fmt"

func longestPalindrome(s string) string {
	r := result005{input: s}
	r.findLongestPalindrome()
	return r.palindrome
}

type result005 struct {
	input      string
	palindrome string
}

func (r *result005) findLongestPalindrome() {
	for i := range r.input {
		words := getPs(r.input, i)
		//fmt.Println(fmt.Sprintf("(i, words): (%d,%v)", i, words))
		for _, w := range words {
			if len(w) > len(r.palindrome) {
				r.palindrome = w
			}
		}
	}
}

func getPs(s string, i int) []string {
	result := []string{string(s[i])}

	pUint8 := s[i]
	p := string(pUint8)
	for j := 1; j < len(s)-i; j++ {
		left := i - j
		right := i + j
		if left >= 0 {
			if s[left] == s[right] {
				p = fmt.Sprintf("%s%s%s", string(s[left]), p, string(s[right]))
				result = append(result, p)
			} else {
				break
			}
		}
	}

	p = ""
	if i+1 < len(s) {
		for j := 0; j < len(s)-i-1; j++ {
			left := i - j
			right := i + 1 + j
			if left >= 0 {
				if s[left] == s[right] {
					p = fmt.Sprintf("%s%s%s", string(s[left]), p, string(s[right]))
					result = append(result, p)
				} else {
					break
				}
			}
		}
	}

	return result
}
