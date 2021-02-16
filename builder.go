package sqlf

// Builder is responsible to build SelectStatements with a default configuration.
type Builder interface {
	Placeholder(format PlaceholderFormat) Builder
	Select(fields ...string) Select
	Insert(tableName string, fields ...interface{}) Insert
	Delete(tableName ...string) Delete
	Update(tableName ...string) Update
}
