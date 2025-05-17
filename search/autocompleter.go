package search

type Autocompleter[T any] struct {
	idx TagIndex[T]
}

func NewAutocompleter[T any](idx TagIndex[T]) *Autocompleter[T] {
	return &Autocompleter[T]{idx: idx}
}

func (a Autocompleter[_]) Autocomplete(tags []string, limit int) ([][]TagInfo, error) {

	// ???

	panic("TODO")
}
