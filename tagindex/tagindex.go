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

type TagIndex interface {
	Query(tags []string, limit int) QueryResult
	Put(e ...Entity)
}
