package sqlf

import "strings"

var (
	sqlDeleteStatement        = []byte("DELETE FROM ")
	sqlDeleteCascadeStatement = []byte("DELETE CASCADE FROM ")
)

type DeleteStatement struct {
	placeholderFormat PlaceholderFormat
	cascade           bool
	from              string
	where             []Sqlizer
	suffix            string
}

// Placeholder defines the placeholder format that should be used for this delete statement.
func (d *DeleteStatement) Placeholder(placeholder PlaceholderFormat) Delete {
	d.placeholderFormat = placeholder
	return d
}

// Cascade enables the CASCADE option.
func (d *DeleteStatement) Cascade() Delete {
	d.cascade = true
	return d
}

// From defines what table will be deleted.
func (d *DeleteStatement) From(tableName string) Delete {
	d.from = tableName
	return d
}

// Where appends a condition. If called multiples, the conditions will be appended.
//
// The conditions added will use the AND operator.
func (d *DeleteStatement) Where(condition string, args ...interface{}) Delete {
	if d.where == nil {
		d.where = []Sqlizer{
			Condition(condition, args...),
		}
		return d
	}
	d.where = append(d.where, Condition(condition, args...))
	return d
}

// WhereClause appends any Sqlizer to serve as where.
//
// The conditions added will use the AND operator.
func (d *DeleteStatement) WhereClause(conditions ...Sqlizer) Delete {
	if d.where == nil {
		d.where = conditions
		return d
	}
	d.where = append(d.where, conditions...)
	return d
}

// Suffix adds a suffix to the DELETE statement. That can be useful for
// extending the SQL for uncovered database technologies.
func (d *DeleteStatement) Suffix(suffix string) Delete {
	d.suffix = suffix
	return d
}

// ToSQL generates the SQL and returns it, alongside its params.
func (d *DeleteStatement) ToSQL() (string, []interface{}, error) {
	sb := new(strings.Builder)
	args := make([]interface{}, 0)
	err := d.ToSQLFast(sb, &args)
	if err != nil {
		return "", nil, err
	}
	if d.placeholderFormat != nil {
		sql, err := d.placeholderFormat.Replace(sb.String())
		if err != nil {
			return "", nil, err
		}
		return sql, args, nil
	}
	return sb.String(), args, nil
}

// ToSQLFast generates the SQL and returns it, alongside its params.
func (d *DeleteStatement) ToSQLFast(sb *strings.Builder, args *[]interface{}) error {
	if d.cascade {
		sb.Write(sqlDeleteCascadeStatement)
	} else {
		sb.Write(sqlDeleteStatement)
	}
	sb.WriteString(d.from)
	if len(d.where) > 0 {
		sb.Write(sqlWhereClause)
		for idx, condition := range d.where {
			if idx > 0 {
				sb.Write(sqlConditionAnd)
			}
			err := RenderInterfaceAsSQL(sb, args, condition)
			if err != nil {
				return err
			}
		}
	}

	if d.suffix != "" {
		sb.Write(sqlSpace)
		sb.WriteString(d.suffix)
	}

	return nil
}
