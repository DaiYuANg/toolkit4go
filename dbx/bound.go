package dbx

type BoundQuery struct {
	SQL  string
	Args []any
}
