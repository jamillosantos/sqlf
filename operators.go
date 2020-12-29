package sqlf

import "strings"

var (
	sqlOperatorNot          = []byte("NOT ")
	sqlOperatorBracketOpen  = []byte("(")
	sqlOperatorBracketClose = []byte(")")
	sqlOperatorAnd          = []byte(" AND ")
	sqlOperatorOr           = []byte(" OR ")
)

// operator is a generic operator helper that will render a SQL by joining all
// input parts with `space`.
type operator struct {
	separator []byte
	parts     []interface{}
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (c *operator) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.Write(sqlOperatorBracketOpen)
	for idx, part := range c.parts {
		if idx > 0 {
			sb.Write(c.separator)
		}
		err := RenderInterfaceAsSQL(sb, args, part)
		if err != nil {
			return err
		}
	}
	sb.Write(sqlOperatorBracketClose)
	return nil
}

// And receive many conditions and generates their SQL with the AND operator
// between the conditions.
func And(conditions ...interface{}) FastSqlizer {
	return &operator{
		separator: sqlOperatorAnd,
		parts:     conditions,
	}
}

// Or receive many conditions and generates their SQL with the OR operator
// between the conditions.
func Or(conditions ...interface{}) FastSqlizer {
	return &operator{
		separator: sqlOperatorOr,
		parts:     conditions,
	}
}

type notOperator struct {
	condition FastSqlizer
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (not *notOperator) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.Write(sqlOperatorNot)
	return not.condition.ToSQLFast(sb, args)
}

// Not negates whatever conditions are passed returning a rendered SQL.
func Not(condition FastSqlizer) FastSqlizer {
	return &notOperator{
		condition: condition,
	}
}
