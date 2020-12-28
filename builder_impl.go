package sqlf

type builder struct {
	placeholder PlaceholderFormat
}

// NewBuilder returns a new instance of the default implementation of the `Builder`.
func NewBuilder() Builder {
	return &builder{}
}

func (b *builder) Placeholder(format PlaceholderFormat) Builder {
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

func (b *builder) Delete(tableName string) Delete {
	return &DeleteStatement{
		placeholderFormat: b.placeholder,
		from:              tableName,
	}
}
