package search

type Autocompleter struct {
	idx TagIndex
}

func NewAutocompleter(idx TagIndex) *Autocompleter {
	return &Autocompleter{idx: idx}
}

func (a Autocompleter) Autocomplete(tags []string, limit int) ([][]TagInfo, error) {

	// ???

	panic("TODO")
}
