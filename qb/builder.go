package qb

type QueryBuilder interface {
	// ToCQL returns the CQL query string from operations performed on the Builder.
	ToCQL() string

	// Build returns the CQL query string and any error encountered during the build process.
	// A set of validation rules will be performed on the query during the build process,
	// and any errors encountered will be returned.
	Build() (string, error)

	// QueryValues returns the values to be used in the query, based on values provided to the builder.
	QueryValues() []any
}
