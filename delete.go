package sqlf

// Delete describes how a DELETE will behave into the sqlf.
type Delete interface {
	Sqlizer
	FastSqlizer

	// Placeholder defines the placeholder format that should be used for this delete statement.
	Placeholder(placeholder PlaceholderFormat) Delete

	// Cascade enables the CASCADE option.
	Cascade() Delete

	// From defines what table will be deleted.
	From(tableName string) Delete

	// Where appends a condition. If called multiples, the conditions will be appended.
	//
	// The conditions added will use the AND operator.
	Where(condition string, args ...interface{}) Delete

	// WhereClause appends any Sqlizer to serve as where.
	//
	// The conditions added will use the AND operator.
	WhereClause(conditions ...FastSqlizer) Delete

	// Suffix adds a suffix to the DELETE statement. That can be useful for
	// extending the SQL for uncovered database technologies.
	Suffix(suffix string) Delete
}
