package lc

import (
	"fmt"
)

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	var tree *node
	tree = insert(nums1, nums2, tree)

	if d(tree.l) == d(tree.r) {
		return float64(tree.v)
	}
	return float64((float64(tree.v + tree.l.v)) / 2)
}

func insert(nums1 []int, nums2 []int, tree *node) *node {

	if len(nums1) == 0 && len(nums2) == 0 {
		return tree
	}
	v, flag := getVAndFlag(nums1, nums2)
	newNode := &node{v, nil, nil}

	if flag {
		nums1 = nums1[1:]
	} else {
		nums2 = nums2[1:]
	}

	if tree == nil {
		tree = newNode
		return insert(nums1, nums2, tree)
	}

	if tree.l == nil {
		newNode.l = tree
		tree = newNode
		return insert(nums1, nums2, tree)
	}

	p := tree
	for p.r != nil {
		p = p.r
	}
	p.r = newNode

	if d(tree.l) == d(tree.r) {
		return insert(nums1, nums2, tree)
	}

	if d(tree.l) < d(tree.r) {
		r := tree.r
		tree.r = nil
		r.l = tree
		tree = r
		return insert(nums1, nums2, tree)
	}

	panic(fmt.Sprintf("d(tree.l) < d(tree.r) occurred: %d, %d", d(tree.l), d(tree.r)))
}
func getVAndFlag(nums1 []int, nums2 []int) (int, bool) {
	var v int
	var flag bool
	if len(nums1) != 0 && len(nums2) == 0 {
		v = nums1[0]
		flag = true
	}
	if len(nums1) == 0 && len(nums2) != 0 {
		v = nums2[0]
	}
	if len(nums1) != 0 && len(nums2) != 0 {
		v = nums1[0]
		flag = true
		if nums2[0] < nums1[0] {
			v = nums2[0]
			flag = false
		}
	}

	return v, flag
}
