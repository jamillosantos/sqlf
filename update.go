package sqlf

type UpdateFields interface {
	Set(field interface{}, value interface{})
	// TODO(Jota): Add where.
}

type Update interface {
	UpdateFields
}
