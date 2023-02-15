package search

type Operator int

const (
	EQ Operator = iota + 1
	NE
	LIKE
	GT
	LT
	GTE
	LTE
	IN
	NI
	IsNull
	Exists
)

type Filter struct {
	FieldName string
	Value     interface{}
	Operator  Operator
}
