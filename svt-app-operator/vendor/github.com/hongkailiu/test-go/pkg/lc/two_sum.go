package lc

func twoSum(nums []int, target int) []int {
	for i, n := range nums {
		for j, m := range nums {
			if i != j && n+m == target {
				return []int{i, j}
			}
		}
	}
	return nil
}
