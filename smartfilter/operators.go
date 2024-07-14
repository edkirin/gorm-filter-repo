package smartfilter

type Operator string

const (
	OperatorEQ     Operator = "EQ"
	OperatorNE     Operator = "NE"
	OperatorGT     Operator = "GT"
	OperatorGE     Operator = "GE"
	OperatorLT     Operator = "LT"
	OperatorLE     Operator = "LE"
	OperatorLIKE   Operator = "LIKE"
	OperatorILIKE  Operator = "ILIKE"
	OperatorIN     Operator = "IN"
	OperatorNOT_IN Operator = "NOT_IN"
)

var OPERATORS = []Operator{
	OperatorEQ, OperatorNE,
	OperatorGT, OperatorGE, OperatorLT, OperatorLE,
	OperatorLIKE, OperatorILIKE,
	OperatorIN, OperatorNOT_IN,
}
