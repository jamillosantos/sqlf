package sqlf

// Sqlizer define anything that outputs a SQL.
type Sqlizer interface {
	// ToSql generates the SQL and returns it, alongside its params.
	ToSql() (string, []interface{}, error)
}
