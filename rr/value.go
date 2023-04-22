package rr

// TODO: []byte
// TODO: null

// TODO: list (slice)
// TODO: map
// TODO: int set
// TODO: string set

type V interface {
	Int() int64
	Ints() []int64
	Float() float64
	Bool() bool
	Str() string
	Strs() []string
	Bytes() []byte

	// ???
}

/*
type V struct {
	s  string
	ss []string
	//m   map[string]V
	err error
}

func Str(s string) V {
	return V{s: s}
}

func Strs(s ...string) V {
	return V{ss: s}
}

func Int(n int64) V {
	return V{s: strconv.FormatInt(n, 10)}
}

func Ints(n ...int64) V {
	ss := make([]string, len(n))
	for _, n := range n {
		ss = append(ss, strconv.FormatInt(n, 10))
	}
	return V{ss: ss}
}

func Float(n float64) V {
	return V{s: strconv.FormatFloat(n, 'f', -1, 64)}
}

func Bool(b bool) V {
	return V{s: strconv.FormatBool(b)}
}

func (v V) Str() string {
	return v.s
}

func (v *V) Int() (n int64) {
	n, v.err = strconv.ParseInt(v.s, 10, 64)
	return
}

func (v *V) Float() (n float64) {
	n, v.err = strconv.ParseFloat(v.s, 64)
	return
}

func (v *V) Bool() (b bool) {
	b, v.err = strconv.ParseBool(v.s)
	return
}

func (v V) Err() error {
	return v.err
}
*/
