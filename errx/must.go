package errx

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must1 panics if err!=nil, otherwise returns v.
func Must1[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Must2 panics if err!=nil, otherwise returns v1 and v2.
func Must2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// Must3 panics if err!=nil, otherwise returns v1, v2, and v3.
func Must3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3
}

// Must4 panics if err!=nil, otherwise returns v1, v2, v3, and v4.
func Must4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4
}
