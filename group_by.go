package sqlf

type GroupByClause struct {
	fields []interface{}
	having []FastSqlizer
}

// Fields defines the fields that the SQL GROUP BY will group.
func (groupBy *GroupByClause) Fields(fields ...interface{}) GroupBy {
	groupBy.fields = fields
	return groupBy
}

// Having defines the SQL HAVING clause.
func (groupBy *GroupByClause) Having(condition string, args ...interface{}) GroupBy {
	groupBy.having = []FastSqlizer{Condition(condition, args...)}
	return groupBy
}

// HavingClause defines the SQL HAVING clause.
func (groupBy *GroupByClause) HavingClause(criteria ...FastSqlizer) GroupBy {
	groupBy.having = criteria
	return groupBy
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (groupBy *GroupByClause) ToSQLFast(sb SQLWriter, args *[]interface{}) error {
	sb.Write(sqlSelectGroupByClause)
	for idx, field := range groupBy.fields {
		if idx > 0 {
			sb.Write(sqlComma)
		}
		err := RenderInterfaceAsSQL(sb, args, field)
		if err != nil {
			return err
		}
	}
	if len(groupBy.having) > 0 {
		sb.Write(sqlSelectHavingClause)
		for idx, condition := range groupBy.having {
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
