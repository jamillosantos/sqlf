package sqlf

// Update describes how a UPDATE will behave into the sqlf.
type Update interface {
	Sqlizer
	FastSqlizer

	// Placeholder defines the placeholder format that should be used for this update statement.
	Placeholder(placeholder PlaceholderFormatFactory) Update

	// Table defines what table will be updated.
	Table(tableName ...string) Update

	// Set define what fields will be updated, alongside its values.
	//
	// The arguments are mixed with the values, in a alternating order. So, the
	// first argument should be the name of the field that will be updated, the second
	// will be its value. The third the name of the second field, the fourth its value
	// and so on. Hence, the number of arguments passed should, always, be even.
	Set(fieldAndValues ...interface{}) Update

	// Where appends a condition. If called multiples, the conditions will be appended.
	//
	// The conditions added will use the AND operator.
	Where(condition string, args ...interface{}) Update

	// WhereClause appends any Sqlizer to serve as where.
	//
	// The conditions added will use the AND operator.
	WhereClause(conditions ...FastSqlizer) Update
}
