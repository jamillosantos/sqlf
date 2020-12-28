package sqlf

import "strings"

// Sqlizer define anything that outputs a SQL.
type Sqlizer interface {
	// ToSQL generates the SQL and returns it, alongside its params.
	ToSQL() (string, []interface{}, error)

	// ToSQLFast generates the SQL and returns it, alongside its params.
	ToSQLFast(*strings.Builder, *[]interface{}) error
}
