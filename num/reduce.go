package num

import "golang.org/x/exp/constraints"

func Sum[N Number](n ...N) (sum N) {
	for i := range n {
		sum += n[i]
	}
	return
}

func Min[T constraints.Ordered](n ...T) (min T) {
	if len(n) == 0 {
		return
	}
	min = n[0]
	for i := range n {
		if n[i] < min {
			min = n[i]
		}
	}
	return
}

func Max[T constraints.Ordered](n ...T) (max T) {
	if len(n) == 0 {
		return
	}
	max = n[0]
	for i := range n {
		if n[i] > max {
			max = n[i]
		}
	}
	return
}
