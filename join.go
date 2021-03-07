package sqlf

// JoinClause is the default implementation for the Join interface.
type JoinClause struct {
	parent   Select
	joinType string
	table    string
	as       string
	on       []FastSqlizer
	using    []interface{}
}

func NewJoinClause(table ...string) Join {
	join := &JoinClause{
		table: table[0],
	}
	if len(table) > 1 {
		join.as = table[1]
	}
	return join
}

// Type defines the type of the Join. Ex: INNER, LEFT, OUTER, etc.
func (join *JoinClause) Type(joinType string) Join {
	join.joinType = joinType
	return join
}

// On define the on criteria.
func (join *JoinClause) On(condition string, params ...interface{}) Select {
	join.on = []FastSqlizer{Condition(condition, params...)}
	return join.parent
}

// OnClause define the on criteria.
func (join *JoinClause) OnClause(criteria ...FastSqlizer) Select {
	join.on = criteria
	return join.parent
}

// Using defines the using directive.
func (join *JoinClause) Using(fields ...interface{}) Select {
	join.using = fields
	return join.parent
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (join *JoinClause) ToSQLFast(sb SQLWriter, args *[]interface{}) error {
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
			err := RenderInterfaceAsSQL(sb, args, field)
			if err != nil {
				return err
			}
		}
		sb.Write(sqlBracketClose)
	}

	return nil
}
