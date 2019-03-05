package lc

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	return addNode(l1, l2, 0)
}

func addNode(l1, l2 *ListNode, i int) *ListNode {
	if l1 == nil && l2 == nil {
		if i == 0 {
			return nil
		}
		return &ListNode{i, nil}
	}

	if l1 == nil && l2 != nil {
		return &ListNode{(l2.Val + i) % 10, addNode(nil, l2.Next, (l2.Val+i)/10)}
	}

	if l1 != nil && l2 == nil {
		return &ListNode{(l1.Val + i) % 10, addNode(l1.Next, nil, (l1.Val+i)/10)}
	}

	return &ListNode{(l1.Val + l2.Val + i) % 10, addNode(l1.Next, l2.Next, (l1.Val+l2.Val+i)/10)}
}
