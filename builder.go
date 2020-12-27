package sqlf

type Builder interface {
	Placeholder(format PlaceholderFormat) Builder
	Select(fields ...string) Select
}
