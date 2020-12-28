package sqlf

import "strings"

type GroupByClause struct {
	fields []interface{}
	having []Sqlizer
}

// Fields defines the fields that the SQL GROUP BY will group.
func (groupBy *GroupByClause) Fields(fields ...interface{}) GroupBy {
	groupBy.fields = fields
	return groupBy
}

// Having defines the SQL HAVING clause.
func (groupBy *GroupByClause) Having(condition string, args ...interface{}) GroupBy {
	groupBy.having = []Sqlizer{Condition(condition, args...)}
	return groupBy
}

// HavingClause defines the SQL HAVING clause.
func (groupBy *GroupByClause) HavingClause(criteria ...Sqlizer) GroupBy {
	groupBy.having = criteria
	return groupBy
}

// ToSQL generates the SQL and returns it, alongside its params.
func (groupBy *GroupByClause) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0, 2)
	err := groupBy.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (groupBy *GroupByClause) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
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
