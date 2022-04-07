package tagindex

type Entity struct {
	ID    string
	Order int64
	Tags  []string
}

type QueryResult struct {
	Data       []Entity
	TotalCount int
}

type TagInfo struct {
	Tag   string
	Count int
}

type TagIndex interface {
	Query(tags []string, limit int) (QueryResult, error)
	Put(e ...Entity) error
	Commit() error
	Tags(prefix string, limit int) ([]TagInfo, error)
}
