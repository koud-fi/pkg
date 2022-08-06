package rx

func Random(min, max int) Iter[N] {
	return FuncIter(func() ([]N, bool, error) {
		return []N{rand.Intn(max - min) + min}, true, nil
	})
}
