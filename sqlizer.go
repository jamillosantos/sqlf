package sqlf

import (
	"fmt"
	"io"
)

// SQLWriter is an interface that abstracts a buffer used to write SQLs. This approach is used to enable
// placeholder replacers to be used by wrapping this `SQLWriter` into another structure. Check `placeholders.go`
// for more information.
type SQLWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
	fmt.Stringer
}

// Sqlizer define anything that outputs a SQL in a simplified form.
type Sqlizer interface {
	// ToSQL generates the SQL and returns it, alongside its params.
	ToSQL() (string, []interface{}, error)
}

// FastSqlizer define anything that outputs a SQL.
type FastSqlizer interface {
	// ToSQLFast generates the SQL and returns it, alongside its params.
	ToSQLFast(SQLWriter, *[]interface{}) error
}
