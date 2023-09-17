package rx

type Lens[T any] interface {
	Get() (T, error)
	Set(T) error
}
