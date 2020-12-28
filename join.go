package sqlf

import (
	"fmt"
	"strings"
)

// JoinClause is the default implementation for the Join interface.
type JoinClause struct {
	parent   Select
	joinType string
	table    string
	as       string
	on       []Sqlizer
	using    []interface{}
}

// Type defines the type of the Join. Ex: INNER, LEFT, OUTER, etc.
func (join *JoinClause) Type(joinType string) Join {
	join.joinType = joinType
	return join
}

// Table defines the table name.
func (join *JoinClause) Table(tableName ...string) Join {
	if len(tableName) > 0 {
		join.table = tableName[0]
	}
	if len(tableName) > 1 {
		join.as = tableName[1]
	}
	return join
}

// As define the table alias.
func (join *JoinClause) As(name string) Join {
	join.as = name
	return join
}

// On define the on criteria.
func (join *JoinClause) On(condition string, params ...interface{}) Select {
	join.on = []Sqlizer{Condition(condition, params...)}
	return join.parent
}

// OnClause define the on criteria.
func (join *JoinClause) OnClause(criteria ...Sqlizer) Select {
	join.on = criteria
	return join.parent
}

// Using defines the using directive.
func (join *JoinClause) Using(fields ...interface{}) Select {
	join.using = fields
	return join.parent
}

// ToSQL generates the SQL and returns it, alongside its params.
func (join *JoinClause) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0, 2)
	err := join.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (join *JoinClause) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.WriteString(join.joinType)
	sb.Write(sqlSelectJoinClause)
	sb.WriteString(join.table)

	// If `as` is not defined, don't append it.
	if join.as != "" {
		sb.Write(sqlSelectAsClause)
		sb.WriteString(join.as)
	}

	// Supposely ON and USING cannot be used together. Let the user deal with it.

	// ON added only if defined.
	if len(join.on) > 0 {
		sb.Write(sqlSelectJoinOnClause)
		for idx, join := range join.on {
			// By default criteria is joined by `AND` condition.
			if idx > 0 {
				sb.Write(sqlConditionAnd)
			}

			//
			err := join.ToSQLFast(sb, args)
			if err != nil {
				return err
			}
		}
	}

	// USING added only if defined.
	if len(join.using) > 0 {
		sb.Write(sqlSelectJoinUsingClause)
		sb.Write(sqlBracketOpen)
		for idx, field := range join.using {
			if idx > 0 {
				sb.Write(sqlComma)
			}
			// Fields are supported as string, fmt.Stringer or `Sqlizer`. This
			// should provide plenty flexibility a wide use cases.
			switch f := field.(type) {
			case string:
				sb.WriteString(f)
			case []byte:
				sb.Write(f)
			case fmt.Stringer:
				sb.WriteString(f.String())
			case Sqlizer:
				err := f.ToSQLFast(sb, args)
				if err != nil {
					return err
				}
			}
		}
		sb.Write(sqlBracketClose)
	}

	return nil
}
