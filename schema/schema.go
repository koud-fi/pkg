package schema

type Schema struct {
	Definitions map[string]Type `json:"definitions,omitempty"`
}
