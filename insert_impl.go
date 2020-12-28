package sqlf

import (
	"errors"
	"strings"
)

var (
	sqlInsertStatement       = []byte("INSERT INTO ")
	sqlInsertValuesClause    = []byte(" VALUES ")
	sqlInsertReturningClause = []byte(" RETURNING ")
)

var (
	// ErrMismatchFieldsAndValuesCount is returned when the length of values is not compatible with the length
	// of fields.
	//
	// This means that the len(values) is not multiple of len(fields). (Check the `Insert.Values` method documentation).
	ErrMismatchFieldsAndValuesCount = errors.New("the amount values is not compatible with the amount of fields")
)

// InsertStatement is the default implementation of the `Insert` interface.
type InsertStatement struct {
	placeholderFormat PlaceholderFormat
	tableName         string
	fields            []interface{}
	values            []interface{}
	selectStatement   Select
	returning         []interface{}
	suffix            string
	suffixArgs        []interface{}
}

// Placeholder defines the placeholder format that should be used for this insert statement.
func (insert *InsertStatement) Placeholder(placeholder PlaceholderFormat) Insert {
	insert.placeholderFormat = placeholder
	return insert
}

// Into defines what table the data will be inserted on. `fields` are the same as `Fields` method.
func (insert *InsertStatement) Into(tableName string, fields ...interface{}) Insert {
	insert.tableName = tableName
	if len(fields) > 0 {
		insert.fields = fields
	}
	return insert
}

// Fields define what fields will be defined into the insert. If there is a set of fields already defined
// this method will replace them.
func (insert *InsertStatement) Fields(fields ...interface{}) Insert {
	insert.fields = fields
	return insert
}

// AddFields will append the fields to the current list.
func (insert *InsertStatement) AddFields(fields ...interface{}) Insert {
	if insert.fields == nil {
		insert.fields = fields
		return insert
	}
	insert.fields = append(insert.fields, fields...)
	return insert
}

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
func (insert *InsertStatement) Values(values ...interface{}) Insert {
	if insert.values == nil {
		insert.values = values
		return insert
	}
	insert.values = append(insert.values, values...)
	return insert
}

// Select defines a select that will be inserted.
//
// Below an example of how this would be used in plain SQL. Ex:
//
//     INSERT INTO users (name, email, password) SELECT name, email, "12345" FROM employees;
//
// Select and Values are not compatible. Whenever there is a select defined, Values should be ignored.
//
func (insert *InsertStatement) Select(callback func(Select)) Insert {
	insert.selectStatement = &SelectStatement{
		placeholderFormat: insert.placeholderFormat,
	}
	callback(insert.selectStatement)
	return insert
}

// Returning defines the RETURNING clause defined for Postgres.
func (insert *InsertStatement) Returning(fields ...interface{}) Insert {
	insert.returning = fields
	return insert
}

// OnConflict defines the ON CONFLICT clause for Postgres.
func (insert *InsertStatement) OnConflict(callback func(InsertConflict)) Insert {
	panic("not implemented") // TODO(Jota): To implement
}

// Suffix defines a suffix that will be appended at the end of the insert clause. This can be used to extend the
// uses (like "ON DUPLICATE KEY" on MySQL) for other database technologies.
func (insert *InsertStatement) Suffix(suffix string, args ...interface{}) Insert {
	insert.suffix = suffix
	insert.suffixArgs = args
	return insert
}

// ToSQL generates the SQL and returns it, alongside its params.
func (insert *InsertStatement) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0)
	err := insert.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	if insert.placeholderFormat != nil {
		sql, err := insert.placeholderFormat.Replace(sb.String())
		if err != nil {
			return "", nil, err
		}
		return sql, args, nil
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (insert *InsertStatement) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	lenFields := len(insert.fields)
	// if the selectStatement is not defined AND if the values count is multiple of the fields count.
	if insert.selectStatement == nil && len(insert.values)%lenFields != 0 {
		return ErrMismatchFieldsAndValuesCount
	}

	// Writing >> INSERT INTO <<
	sb.Write(sqlInsertStatement)

	// Writing insert into >> <TABLENAME> <<
	sb.WriteString(insert.tableName)
	sb.Write(sqlSpace)

	// Writing insert into <tablename> >> (<FIELDS>) <<
	sb.Write(sqlBracketOpen)
	for idx, field := range insert.fields {
		if idx > 0 {
			sb.Write(sqlComma)
		}
		err := RenderInterfaceAsSQL(sb, args, field)
		if err != nil {
			return err
		}
	}
	sb.Write(sqlBracketClose)

	if insert.selectStatement == nil {
		// Writting insert into <tablename> (<fields>) >> VALUES (<VALUES>) <<
		sb.Write(sqlInsertValuesClause)
		recordCount := len(insert.values) / lenFields
		for i := 0; i < recordCount; i++ {
			if i > 0 {
				sb.Write(sqlComma)
			}
			sb.Write(sqlBracketOpen)
			sb.WriteString(Placeholders(lenFields))
			sb.Write(sqlBracketClose)
		}
		*args = append(*args, insert.values...)
	} else {
		// Writting insert into <tablename> (<fields>) >> SELECT ... FROM ... <<
		sb.Write(sqlSpace)
		err := insert.selectStatement.ToSQLFast(sb, args)
		if err != nil {
			return err
		}
	}

	if len(insert.returning) > 0 {
		// Writting insert into ... values (...) >> RETURNING <fields> <<

		sb.Write(sqlInsertReturningClause)
		//
		for idx, field := range insert.returning {
			if idx > 0 {
				sb.Write(sqlComma)
			}
			err := RenderInterfaceAsSQL(sb, args, field)
			if err != nil {
				return err
			}
		}
	}

	if insert.suffix != "" {
		sb.Write(sqlSpace)
		sb.WriteString(insert.suffix)
		if len(insert.suffixArgs) > 0 {
			*args = append(*args, insert.suffixArgs...)
		}
	}
	return nil
}
