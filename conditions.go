package sqlf

import "strings"

var (
	sqlSpace = []byte(" ")
)

type condition struct {
	sql  string
	args []interface{}
}

// ToSQL generates the SQL and returns it, alongside its params.
func (condition *condition) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.WriteString(condition.sql)
	if len(condition.args) > 0 {
		*args = append(*args, condition.args...)
	}
	return nil
}

func Condition(sql string, args ...interface{}) FastSqlizer {
	return &condition{
		sql:  sql,
		args: args,
	}
}
