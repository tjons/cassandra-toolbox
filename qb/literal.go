package qb

type literal struct {
	value       string
	singleQuote bool
}

// Literal is used to insert a raw CQL fragment into a query.
// It will not be escaped or quoted in any way, so it should be used with caution.
func Literal(value string) literal {
	return literal{value: value}
}

// CqlFunction is used to insert a CQL function into a query.
// It will not be escaped or quoted in any way, so it should be used with caution.
func CqlFunction(value string) literal {
	return literal{value: value}
}

// StringLiteral is used to insert a string literal into a query.
// It will be quoted with single quotes.
func StringLiteral(value string) literal {
	return literal{value: value, singleQuote: true}
}
