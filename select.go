package sqlf

// PlaceholderFormat enables the system to use different placeholder formats.
type PlaceholderFormat interface {
	Replace(sql string) (string, error)
}

// Join represents a SQL Join.
type Join interface {
	Sqlizer

	// Type defines the type of the Join. Ex: INNER, LEFT, OUTER, etc.
	Type(joinType string) Join

	// Table defines the table name.
	Table(tableName ...string) Join

	// As define the table alias.
	As(name string) Join

	// On define the on criteria.
	On(criteria ...Sqlizer) Select

	// Using defines the using directive.
	Using(fields ...interface{}) Select
}

// GroupBy represents a SQL GROUP BY clause.
type GroupBy interface {
	Sqlizer

	// Fields defines the fields that the SQL GROUP BY will group.
	Fields(fields ...interface{}) GroupBy

	// Having defines the SQL HAVING clause.
	Having(criteria ...Sqlizer) Select

	// Query returns the Query that created this instance.
	Query() Select
}

// OrderBy represents a SQL GROUP BY clause.
type OrderBy interface {
	Sqlizer

	// Asc adds fields to the SQL ORDER BY clause on an ascending order.
	Asc(fields ...interface{}) OrderBy

	// Desc adds fields to the SQL ORDER BY clause on an descending order.
	Desc(fields ...interface{}) OrderBy

	// Query returns the Query that created this instance.
	Query() Select
}

// Select represents a SQL SELECT statement.
type Select interface {
	Sqlizer

	// Select defines the fields that will be returned by the query.
	Select(fields ...interface{}) Select

	// AddSelect adds fields to the existing list of fields that will be selected.
	AddSelect(fields ...interface{}) Select

	// Distinct enables the SQL SELECT DISTINCT clause.
	Distinct() Select

	// From defines the SQL SELECT FROM clause.
	From(table ...string) Select

	// As defines alias for the table.
	As(tableAlias string) Select

	// JoinClause adds a JOIN to the select.
	JoinClause(joinType string, table string) Join

	// InnerJoin adds a INNER JOIN to the select.
	InnerJoin(table string) Join

	// OuterJoin adds a OUTER JOIN to the select.
	OuterJoin(table string) Join

	// LeftJoin adds a LEFT JOIN to the select.
	LeftJoin(table string) Join

	// RightJoin adds a LEFT JOIN to the select.
	RightJoin(table string) Join

	// Where adds a criteria for the
	Where(criteria ...Sqlizer) Select

	// GroupBy adds a SQL GROUP BY clause and returns the Query itself. For more options (like HAVING) use `GroupByX`.
	GroupBy() Select

	// GroupByX adds a SQL GROUP BY clause and returns the GroupBy itself for further configuration.
	GroupByX() GroupBy

	// OrderBy adds a SQL GROUP BY clause and returns the Query itself. For more options (like HAVING) use `OrderByX`.
	OrderBy(fields ...interface{}) OrderBy

	// OrderByX adds a SQL GROUP BY clause and returns the OrderBy itself for further configuration.
	OrderByX() OrderBy

	// Limit defines the SQL LIMIT clause.
	Limit(limits ...interface{}) Select

	// Placeholder defines what placeholder format is going to be used for this query.
	//
	// Usually it will be automatically defined by the `Builder`.
	Placeholder(placeholder PlaceholderFormat) Select
}
