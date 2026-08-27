package repo

import (
	"embed"
)

type QueriesReader struct {
	eq  embed.FS
	err error
}

func NewQueriesReader(eq embed.FS) *QueriesReader {
	return &QueriesReader{eq: eq}
}

func (r *QueriesReader) GetQuery(queryName string) string {
	if r.err != nil {
		return ""
	}
	var file []byte
	file, r.err = r.eq.ReadFile(queryName + ".sql")
	if r.err != nil {
		return ""
	}
	return string(file)
}

func (r *QueriesReader) Err() error {
	return r.err
}
