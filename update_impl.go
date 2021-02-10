package sqlf

import (
	"errors"
	"strings"
)

var (
	sqlUpdateStatement       = []byte("UPDATE ")
	sqlUpdateSetClause       = []byte(" SET ")
	sqlUpdateAssignOperation = []byte(" = ")
)

var (
	ErrUpdateInvalidFieldValuePairCount = errors.New("invalid field and value pair count")
)

type UpdateStatement struct {
	placeholderFormat PlaceholderFormat
	tableName         string
	as                string
	fields            []interface{}
	where             []FastSqlizer
}

// Placeholder defines the placeholder format that should be used for this delete statement.
func (update *UpdateStatement) Placeholder(placeholder PlaceholderFormat) Update {
	update.placeholderFormat = placeholder
	return update
}

// Table defines what table will be deleted.
func (update *UpdateStatement) Table(tableName ...string) Update {
	if len(tableName) > 0 {
		update.tableName = tableName[0]
	}
	if len(tableName) > 1 && tableName[1] != "" {
		update.as = tableName[1]
	}
	return update
}

// Set define what fields will be updated, alongside its values.
//
// The arguments are mixed with the values, in a alternating order. So, the
// first argument should be the name of the field that will be updated, the second
// will be its value. The third the name of the second field, the fourth its value
// and so on. Hence, the number of arguments passed should, always, be even.
func (update *UpdateStatement) Set(fieldAndValues ...interface{}) Update {
	if update.fields == nil {
		update.fields = fieldAndValues
		return update
	}
	update.fields = append(update.fields, fieldAndValues...)
	return update
}

// Where appends a condition. If called multiples, the conditions will be appended.
//
// The conditions added will use the AND operator.
func (update *UpdateStatement) Where(condition string, args ...interface{}) Update {
	if update.where == nil {
		update.where = []FastSqlizer{
			Condition(condition, args...),
		}
		return update
	}
	update.where = append(update.where, Condition(condition, args...))
	return update
}

// WhereClause appends any Sqlizer to serve as where.
//
// The conditions added will use the AND operator.
func (update *UpdateStatement) WhereClause(conditions ...FastSqlizer) Update {
	if update.where == nil {
		update.where = conditions
		return update
	}
	update.where = append(update.where, conditions...)
	return update
}

// ToSQL generates the SQL and returns it, alongside its params.
func (update *UpdateStatement) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0)
	err := update.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	if update.placeholderFormat != nil {
		sql, err := update.placeholderFormat.Replace(sb.String())
		if err != nil {
			return "", nil, err
		}
		return sql, args, nil
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (update *UpdateStatement) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	// Writing >> UPDATE <TABLE> SET <<
	sb.Write(sqlUpdateStatement)
	sb.WriteString(update.tableName)
	if update.as != "" {
		sb.Write(sqlSelectAsClause)
		sb.WriteString(update.as)
	}
	sb.Write(sqlUpdateSetClause)

	// Enforce the key-pair for the set clause.
	lenFields := len(update.fields)
	if lenFields%2 != 0 {
		return ErrUpdateInvalidFieldValuePairCount
	}

	// Writing update <table> set >> field = value <<
	for i := 0; i < lenFields; i += 2 {
		if i > 0 {
			sb.Write(sqlComma)
		}
		err := RenderInterfaceAsSQL(sb, args, update.fields[i])
		if err != nil {
			return err
		}
		sb.Write(sqlUpdateAssignOperation)
		err = RenderInterfaceAsArg(sb, args, update.fields[i+1])
		if err != nil {
			return err
		}
	}

	if len(update.where) > 0 {
		// Writing update <table> set field = value >> WHERE <<
		sb.Write(sqlWhereClause)

		// Writing update <table> set field = value where >> <CONDITIONS> <<
		for idx, condition := range update.where {
			if idx > 0 {
				sb.Write(sqlConditionAnd)
			}
			err := RenderInterfaceAsSQL(sb, args, condition)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
