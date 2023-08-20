package search

import "strconv"

const cursorBase = 36

type Entry struct {
	ID    string
	Order int64
	Tags  []string
}

type QueryResult struct {
	Data       []Entry
	TotalCount int
}

func (qr QueryResult) Page(cursor string, limit int) ([]Entry, string) {
	var (
		dataLen    = int64(len(qr.Data))
		start, _   = strconv.ParseInt(cursor, cursorBase, 64)
		end        int64
		nextCursor string
	)
	if limit <= 0 {
		end = dataLen
	} else {
		end = start + int64(limit)
		if end >= dataLen {
			end = dataLen
		} else {
			nextCursor = strconv.FormatInt(end, cursorBase)
		}
	}
	return qr.Data[start:end], nextCursor
}

type TagInfo struct {
	Tag   string
	Count int
}

type TagIndex interface {
	Query(tags []string, limit int) (QueryResult, error)
	Put(e ...Entry)
	Commit() error
	Tags(prefix string) ([]TagInfo, error)
}
