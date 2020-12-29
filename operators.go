package sqlf

import "strings"

var (
	operatorNot          = []byte("NOT ")
	operatorBracketOpen  = []byte("(")
	operatorBracketClose = []byte(")")
	operatorAnd          = []byte(" AND ")
	operatorOr           = []byte(" OR ")
)

// operator is a generic operator helper that will render a SQL by joining all
// input parts with `space`.
type operator struct {
	separator []byte
	parts     []interface{}
}

// ToSQL generates the SQL and returns it, alongside its params.
func (c *operator) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0, 2)
	err := c.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (c *operator) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.Write(operatorBracketOpen)
	for idx, part := range c.parts {
		if idx > 0 {
			sb.Write(c.separator)
		}
		err := RenderInterfaceAsSQL(sb, args, part)
		if err != nil {
			return err
		}
	}
	sb.Write(operatorBracketClose)
	return nil
}

// And receive many conditions and generates their SQL with the AND operator
// between the conditions.
func And(conditions ...interface{}) Sqlizer {
	return &operator{
		separator: operatorAnd,
		parts:     conditions,
	}
}

// Or receive many conditions and generates their SQL with the OR operator
// between the conditions.
func Or(conditions ...interface{}) Sqlizer {
	return &operator{
		separator: operatorOr,
		parts:     conditions,
	}
}

// Not negates whatever conditions are passed returning a rendered SQL.
func Not(conditions ...interface{}) Sqlizer {
	c := make([]interface{}, 2, len(conditions)+3)
	c[0] = operatorNot
	c[1] = operatorBracketOpen
	c = append(c, conditions...)
	c = append(c, operatorBracketClose)
	return &operator{
		separator: operatorAnd,
		parts:     c,
	}
}
