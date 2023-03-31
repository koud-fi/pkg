package rr

// TODO: []byte
// TODO: bool
// TODO: []byte set
// TODO: list (slice)
// TODO: map
// TODO: number
// TODO: number set
// TODO: null
// TODO: string
// TODO: string set

type V struct {
	data any
}

func Value(data any) V {

	// TODO: normalization

	return V{data}
}

func (v V) Any() any { return v.data }
