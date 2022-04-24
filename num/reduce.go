package num

func Sum(n ...int) int {
	sum := 0
	for i := range n {
		sum += n[i]
	}
	return sum
}

func Min(n ...int) int {
	if len(n) == 0 {
		return 0
	}
	min := n[0]
	for i := range n {
		if n[i] < min {
			min = n[i]
		}
	}
	return min
}

func Max(n ...int) int {
	if len(n) == 0 {
		return 0
	}
	max := n[0]
	for i := range n {
		if n[i] > max {
			max = n[i]
		}
	}
	return max
}
