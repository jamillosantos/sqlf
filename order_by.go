package sqlf

import "strings"

type OrderByClause struct {
	fields []interface{}
}

type orderByDesc struct {
	value interface{}
}

// ToSQL generates the SQL and returns it, alongside its params.
func (desc *orderByDesc) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0)
	err := desc.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	return sb.String(), args, nil
}

// ToSQL generates the SQL and returns it, alongside its params.
func (desc *orderByDesc) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	err := RenderInterfaceAsSQL(sb, args, desc.value)
	if err != nil {
		return err
	}
	sb.Write(sqlSelectOrderByDescClause)
	return nil
}

// Asc adds fields to the SQL ORDER BY clause on an ascending order.
func (orderBy *OrderByClause) Asc(fields ...interface{}) OrderBy {
	if orderBy.fields == nil {
		orderBy.fields = fields
	} else {
		orderBy.fields = append(orderBy.fields, fields...)
	}
	return orderBy
}

// Desc adds fields to the SQL ORDER BY clause on an descending order.
func (orderBy *OrderByClause) Desc(fields ...interface{}) OrderBy {
	if orderBy.fields == nil {
		orderBy.fields = make([]interface{}, 0, len(fields))
	}
	for _, field := range fields {
		orderBy.fields = append(orderBy.fields, &orderByDesc{
			value: field,
		})
	}
	return orderBy
}

// ToSQL generates the SQL and returns it, alongside its params.
func (orderBy *OrderByClause) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0)
	err := orderBy.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (orderBy *OrderByClause) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	sb.Write(sqlSelectOrderByClause)
	for idx, field := range orderBy.fields {
		if idx > 0 {
			sb.Write(sqlComma)
		}
		err := RenderInterfaceAsSQL(sb, args, field)
		if err != nil {
			return err
		}
	}
	return nil
}
