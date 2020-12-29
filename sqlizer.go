package sqlf

import "strings"

type Sqlizer interface {
	// ToSQL generates the SQL and returns it, alongside its params.
	ToSQL() (string, []interface{}, error)
}

// FastSqlizer define anything that outputs a SQL.
type FastSqlizer interface {
	// ToSQLFast generates the SQL and returns it, alongside its params.
	ToSQLFast(*strings.Builder, *[]interface{}) error
}
