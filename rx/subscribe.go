package rx

import "log"

func ForEach[T any](it Iter[T], fn func(v T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}

func Log[T any](it Iter[T]) {
	ForEach(it, func(v T) error {
		log.Print(v)
		return nil
	})
}
