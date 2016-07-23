package db

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type JSONRow struct {
	JSONString string `db:"json_string" json:"json_string"`
}

type JSONRowSlice []JSONRow

// JSONConcat concatenates []JSONRow{{"key":"value"},{"key":"value"}, ...}
// into `[{"key":"value"},{"key":"value"},...]`
// Uses Buffer.WriteString() but ignores returned error/count
func (rows JSONRowSlice) concatJSON() string {
	if len(rows) == 0 {
		return "[]"
	}

	defer func() {
		if recover() != nil {
			log.Println("bytes.Buffer.WriteString() panicking")
		}
	}()

	var result bytes.Buffer
	_, _ = result.WriteString("[")
	_, _ = result.WriteString(rows[0].JSONString)
	for i := 1; i < len(rows); i++ {
		_, _ = result.WriteString(",")
		_, _ = result.WriteString(rows[i].JSONString)
	}
	_, _ = result.WriteString("]")
	return result.String()
}

func Get(db *sqlx.DB, query string, params ...interface{}) (string, error) {
	rows := []JSONRow{}
	err := db.Select(&rows, jsonQuery(query), params...)
	slice := JSONRowSlice(rows)
	return slice.concatJSON(), err
}

func jsonQuery(query string) string {
	return fmt.Sprintf(`
SELECT row_to_json(row)::text as json_string FROM(
	%s
) row;
`, strings.TrimRight(query, "\n\t ;"))
}
