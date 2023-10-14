package search

import "strconv"

type Entry interface {
	ID() string
	Tags() []string
}

type OrderedEntry interface {
	Entry
	Order() int64
}

type QueryResult[T any] struct {
	Data       []T
	TotalCount int
}

func (qr QueryResult[T]) Page(cursor string, limit int) ([]T, string) {
	const cursorBase = 36
	var (
		dataLen    = int64(len(qr.Data))
		start, _   = strconv.ParseInt(cursor, cursorBase, 64)
		end        int64
		nextCursor string
	)
	if start < 0 {
		start = 0
	}
	if limit <= 0 {
		end = dataLen
	} else {
		if start > dataLen {
			start = dataLen
		}
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

type TagIndex[T Entry] interface {
	Get(id ...string) ([]T, error)
	Query(tags []string, limit int) (QueryResult[T], error)
	Put(e ...T)
	Commit() error
	Tags(prefix string) ([]TagInfo, error)
}
