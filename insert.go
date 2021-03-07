package sqlf

type InsertUpdate interface {
	Set(fieldsAndValues ...interface{}) InsertUpdate
	Where(condition string, args ...interface{}) InsertUpdate
	WhereClause(conditions ...FastSqlizer) InsertUpdate
}

// InsertConflict describes the conflict statement for insertion.
type InsertConflict interface {
	Target(target interface{}) InsertConflict
	DoNothing() InsertConflict
	Update(callback func(InsertUpdate)) InsertConflict
}

// Insert describes how a insert will behave into the sqlf.
type Insert interface {
	Sqlizer
	FastSqlizer

	// Placeholder defines the placeholder format that should be used for this insert statement.
	Placeholder(placeholder PlaceholderFormatFactory) Insert

	// Into defines what table the data will be inserted on. `fields` are the same as `Fields` method.
	Into(tableName string, fields ...interface{}) Insert

	// Fields define what fields will be defined into the insert. If there is a set of fields already defined
	// this method will replace them.
	Fields(fields ...interface{}) Insert

	// AddFields will append the fields to the current list.
	AddFields(fields ...interface{}) Insert

	// Values define what data will be inserted. Values will always append the information.
	//
	// If you are inserting multiple records at once, just call Values as many times you want to. The only requirement
	// is that the len of the total values should be multiple of the amount of fields that were defined.
	//
	// Example:
	//
	//     i := new(InsertStatement)
	//     i.Into("users", "name", "email", "password").Values("Name 1", "email1@email.com", "12345") // Adding one record
	//     i.Values("Name 2", "email2@email.com", "54321") // Adding another record
	//     i.Values("Name 3", "email3@email.com", "13425", "Name 4", "email4@email.com", "52431") // Adding more two records
	//
	Values(values ...interface{}) Insert

	// Select defines a select that will be inserted.
	//
	// Below an example of how this would be used in plain SQL. Ex:
	//
	//     INSERT INTO users (name, email, password) SELECT name, email, "12345" FROM employees;
	//
	// Select and Values are not compatible. Whenever there is a select defined, Values should be ignored.
	//
	Select(callback func(Select)) Insert

	// Returning defines the RETURNING clause defined for Postgres.
	Returning(fields ...interface{}) Insert

	// OnConflict defines the ON CONFLICT clause for Postgres.
	OnConflict(callback func(InsertConflict)) Insert

	// Suffix defines a suffix that will be appended at the end of the insert clause. This can be used to extend the
	// uses (like "ON DUPLICATE KEY" on MySQL) for other database technologies.
	Suffix(suffix string, args ...interface{}) Insert
}
