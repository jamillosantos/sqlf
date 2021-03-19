package sqlf

type builder struct {
	placeholder PlaceholderFormatFactory
}

// NewBuilder returns a new instance of the default implementation of the `Builder`.
func NewBuilder() Builder {
	return &builder{}
}

func (b *builder) Placeholder(format PlaceholderFormatFactory) Builder {
	b.placeholder = format
	return b
}

func (b *builder) Select(fields ...string) Select {
	return &SelectStatement{
		placeholderFormat: b.placeholder,
	}
}

func (b *builder) Insert(into string, fields ...interface{}) Insert {
	return &InsertStatement{
		placeholderFormat: b.placeholder,
		tableName:         into,
		fields:            fields,
	}
}

func (b *builder) Delete(tableName ...string) Delete {
	var t, as string
	if len(tableName) > 0 {
		t = tableName[0]
	}
	if len(tableName) > 1 {
		as = tableName[1]
	}
	return &DeleteStatement{
		placeholderFormat: b.placeholder,
		from:              t,
		as:                as,
	}
}

func (b *builder) Update(tableName ...string) Update {
	var t, as string
	if len(tableName) > 0 {
		t = tableName[0]
	}
	if len(tableName) > 1 {
		as = tableName[1]
	}
	return &UpdateStatement{
		placeholderFormat: b.placeholder,
		tableName:         t,
		as:                as,
	}
}
