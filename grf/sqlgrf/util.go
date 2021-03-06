package sqlgrf

import (
	"database/sql"
	"strconv"
	"strings"
)

func idStr[T ~int32 | ~int64](id []T) string {
	var ids strings.Builder
	for i, id := range id {
		if i > 0 {
			ids.WriteByte(',')
		}
		ids.WriteString(strconv.FormatInt(int64(id), 10))
	}
	return ids.String()
}

func scanRows[T any](rows *sql.Rows, out []T, fn func(*sql.Rows, *T) error) ([]T, error) {
	for rows.Next() {
		var v T
		if err := fn(rows, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Close()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
