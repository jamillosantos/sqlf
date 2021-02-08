package sqlf

import (
	"strings"
)

var (
	sqlComma                   = []byte(", ")
	sqlConditionAnd            = []byte(" AND ")
	sqlSelectClause            = []byte("SELECT ")
	sqlSelectAllFieldsClause   = []byte("*")
	sqlSelectDistinctClause    = []byte("DISTINCT ")
	sqlSelectFromClause        = []byte(" FROM ")
	sqlSelectAsClause          = []byte(" AS ")
	sqlSelectJoinClause        = []byte(" JOIN ")
	sqlSelectJoinOnClause      = []byte(" ON ")
	sqlSelectJoinUsingClause   = []byte(" USING ")
	sqlWhereClause             = []byte(" WHERE ")
	sqlSelectGroupByClause     = []byte(" GROUP BY ")
	sqlSelectHavingClause      = []byte(" HAVING ")
	sqlSelectOrderByClause     = []byte(" ORDER BY ")
	sqlSelectOrderByDescClause = []byte(" DESC")
	sqlSelectLimitClause       = []byte(" LIMIT ")
	sqlSelectOffsetClause      = []byte(" OFFSET ")
	sqlBracketOpen             = []byte("(")
	sqlBracketClose            = []byte(")")
)

type SelectStatement struct {
	table             string
	as                string
	distinct          bool
	fields            []interface{}
	joins             []Join
	where             []FastSqlizer
	groupBy           GroupBy
	orderBy           OrderBy
	limit             interface{}
	offset            interface{}
	placeholderFormat PlaceholderFormat
}

// Select defines the fields that will be returned by the query.
func (s *SelectStatement) Select(fields ...interface{}) Select {
	s.fields = fields
	return s
}

// AddSelect adds fields to the existing list of fields that will be selected.
func (s *SelectStatement) AddSelect(fields ...interface{}) Select {
	if s.fields == nil {
		s.fields = make([]interface{}, 0, len(fields))
	}
	s.fields = append(s.fields, fields...)
	return s
}

// Distinct enables the SQL SELECT DISTINCT clause.
func (s *SelectStatement) Distinct() Select {
	s.distinct = true
	return s
}

// From defines the SQL SELECT FROM clause.
func (s *SelectStatement) From(table ...string) Select {
	if len(table) > 0 {
		s.table = table[0]
	}
	if len(table) > 1 {
		s.as = table[1]
	}
	return s
}

// As defines alias for the table.
func (s *SelectStatement) As(tableAlias string) Select {
	s.as = tableAlias
	return s
}

// JoinClause adds a JOIN to the select.
func (s *SelectStatement) JoinClause(joinType string, tableName ...string) Join {
	join := new(JoinClause)
	join.parent = s
	if joinType != "" {
		join.Type(joinType)
	}
	if len(tableName) > 0 {
		join.table = tableName[0]
	}
	if len(tableName) > 1 {
		join.as = tableName[1]
	}
	if s.joins == nil {
		s.joins = make([]Join, 0, 1)
	}
	s.joins = append(s.joins, join)
	return join
}

// InnerJoin adds a INNER JOIN to the select.
func (s *SelectStatement) InnerJoin(tableName ...string) Join {
	return s.JoinClause("INNER", tableName...)
}

// OuterJoin adds a OUTER JOIN to the select.
func (s *SelectStatement) OuterJoin(tableName ...string) Join {
	return s.JoinClause("OUTER", tableName...)
}

// LeftJoin adds a LEFT JOIN to the select.
func (s *SelectStatement) LeftJoin(tableName ...string) Join {
	return s.JoinClause("LEFT", tableName...)
}

// RightJoin adds a LEFT JOIN to the select.
func (s *SelectStatement) RightJoin(tableName ...string) Join {
	return s.JoinClause("RIGHT", tableName...)
}

// Where adds a criteria for the select.
func (s *SelectStatement) Where(condition string, args ...interface{}) Select {
	if s.where == nil {
		s.where = make([]FastSqlizer, 0, 1)
	}
	s.where = append(s.where, Condition(condition, args...))
	return s
}

// WhereCriteria adds a criteria for the select.
func (s *SelectStatement) WhereCriteria(criteria ...FastSqlizer) Select {
	if s.where == nil {
		s.where = make([]FastSqlizer, 0, len(criteria))
	}
	s.where = append(s.where, criteria...)
	return s
}

// GroupBy adds a SQL GROUP BY clause and returns the Query itself. For more options (like HAVING) use `GroupByX`.
func (s *SelectStatement) GroupBy(fields ...interface{}) Select {
	if len(fields) == 0 {
		s.groupBy = nil
		return s
	}
	s.groupBy = &GroupByClause{
		fields: fields,
	}
	return s
}

// GroupByX adds a SQL GROUP BY clause and returns the GroupBy itself for further configuration.
func (s *SelectStatement) GroupByX(callback func(GroupBy)) Select {
	s.groupBy = &GroupByClause{}
	callback(s.groupBy)
	return s
}

// OrderBy adds a SQL GROUP BY clause and returns the Query itself. For more options (like HAVING) use `OrderByX`.
func (s *SelectStatement) OrderBy(fields ...interface{}) Select {
	if s.orderBy == nil {
		s.orderBy = &OrderByClause{}
	}
	s.orderBy.Asc(fields...)
	return s
}

// OrderByX adds a SQL GROUP BY clause and returns the OrderBy itself for further configuration.
func (s *SelectStatement) OrderByX(callback func(orderBy OrderBy)) Select {
	if s.orderBy == nil {
		s.orderBy = &OrderByClause{}
	}
	callback(s.orderBy)
	return s
}

// Limit defines the SQL LIMIT clause.
func (s *SelectStatement) Limit(limits ...interface{}) Select {
	if len(limits) > 1 {
		s.offset = limits[0]
		s.limit = limits[1]
	} else if len(limits) > 0 {
		s.limit = limits[0]
	}
	return s
}

// Offset defines the SQL OFFSET clause.
func (s *SelectStatement) Offset(offset interface{}) Select {
	s.offset = offset
	return s
}

// Placeholder defines what placeholder format is going to be used for this query.
//
// Usually it will be automatically defined by the `Builder`.
func (s *SelectStatement) Placeholder(placeholder PlaceholderFormat) Select {
	s.placeholderFormat = placeholder
	return s
}

// ToSQL generates the SQL and returns it, alongside its params.
func (s *SelectStatement) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0)
	err := s.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	if s.placeholderFormat != nil {
		sql, err := s.placeholderFormat.Replace(sb.String())
		if err != nil {
			return "", nil, err
		}
		return sql, args, nil
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (s *SelectStatement) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.Write(sqlSelectClause)
	if s.distinct {
		sb.Write(sqlSelectDistinctClause)
	}
	if len(s.fields) == 0 {
		sb.Write(sqlSelectAllFieldsClause)
	} else {
		for idx, field := range s.fields {
			if idx > 0 {
				sb.Write(sqlComma)
			}
			err := RenderInterfaceAsSQL(sb, args, field)
			if err != nil {
				return err
			}
		}
	}
	sb.Write(sqlSelectFromClause)
	sb.WriteString(s.table)
	if s.as != "" {
		sb.Write(sqlSelectAsClause)
		sb.WriteString(s.as)
	}

	for _, join := range s.joins {
		sb.Write(sqlSpace)
		err := join.ToSQLFast(sb, args)
		if err != nil {
			return err
		}
	}

	if len(s.where) > 0 {
		sb.Write(sqlWhereClause)
		for idx, condition := range s.where {
			if idx > 0 {
				sb.Write(sqlConditionAnd)
			}
			err := RenderInterfaceAsSQL(sb, args, condition)
			if err != nil {
				return err
			}
		}
	}

	if s.groupBy != nil {
		err := s.groupBy.ToSQLFast(sb, args)
		if err != nil {
			return err
		}
	}

	if s.orderBy != nil {
		err := s.orderBy.ToSQLFast(sb, args)
		if err != nil {
			return err
		}
	}

	if s.limit != nil {
		sb.Write(sqlSelectLimitClause)
		err := RenderInterfaceAsArg(sb, args, s.limit)
		if err != nil {
			return err
		}
	}

	if s.offset != nil {
		sb.Write(sqlSelectOffsetClause)
		err := RenderInterfaceAsArg(sb, args, s.offset)
		if err != nil {
			return err
		}
	}

	return nil
}
