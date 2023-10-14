package search

type Autocompleter[T Entry] struct {
	idx TagIndex[T]
}

func NewAutocompleter[T Entry](idx TagIndex[T]) *Autocompleter[T] {
	return &Autocompleter[T]{idx: idx}
}

func (a Autocompleter[_]) Autocomplete(tags []string, limit int) ([][]TagInfo, error) {

	// ???

	panic("TODO")
}
