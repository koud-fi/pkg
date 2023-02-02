package rx

type Lens[T any] interface {
	Get() (T, error)
	Set(T) error
}

type valueLens[T any] struct{ v T }

func (vl valueLens[T]) Get() (T, error) { return vl.v, nil }
func (vl *valueLens[T]) Set(v T) error  { vl.v = v; return nil }

func Value[T any](v T) Lens[T] { return &valueLens[T]{v} }
