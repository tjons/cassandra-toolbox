package qb

import (
	"strconv"
	"strings"
)

// DeleteBuilder builds DELETE queries.
type DeleteBuilder interface {
	QueryBuilder

	// Column specifies a column to delete.
	Column(column string) DeleteBuilder

	// From specifies the table to delete from.
	From(table string) DeleteBuilder

	// Using specifies an optional USING clause for the DELETE query.
	// This can be used to specify a USING TIMESTAMP clause.
	// Multiple USING clauses can be specified, and they will be joined
	// with AND in the resulting query.
	// They will be added to the query in the order they were specified.
	Using(updateParam) DeleteBuilder

	// Where specifies a condition for the DELETE query. Multiple calls to
	// Where will be joined with AND in the resulting query.
	Where(column string, ft filterTerm) DeleteBuilder

	// IfExists specifies that the DELETE query should include an IF EXISTS clause.
	IfExists() DeleteBuilder
}

type deleteBuilder struct {
	columns     []string
	table       string
	using       []updateParam
	filterTerms []*filterTerm
	ifExists    bool
	queryValues []any
}

// NewDelete creates a new DeleteBuilder.
func NewDelete() DeleteBuilder {
	return &deleteBuilder{}
}

// Column specifies a column to delete.
func (b *deleteBuilder) Column(column string) DeleteBuilder {
	b.columns = append(b.columns, column)

	return b
}

// Using specifies an optional USING clause for the DELETE query.
// This can be used to specify a USING TIMESTAMP clause.
// Multiple USING clauses can be specified, and they will be joined
// with AND in the resulting query.
// They will be added to the query in the order they were specified.
func (b *deleteBuilder) Using(param updateParam) DeleteBuilder {
	b.using = append(b.using, param)
	return b
}

// IfExists specifies that the DELETE query should include an IF EXISTS clause.
func (b *deleteBuilder) IfExists() DeleteBuilder {
	b.ifExists = true

	return b
}

// From specifies the table to delete from.
func (b *deleteBuilder) From(table string) DeleteBuilder {
	b.table = table

	return b
}

// Where specifies a condition for the DELETE query.
// Multiple calls to Where will be joined with AND in the resulting query.
func (b *deleteBuilder) Where(column string, ft filterTerm) DeleteBuilder {
	ft.column = column
	b.filterTerms = append(b.filterTerms, &ft)

	return b
}

// Build builds the CQL query string for the delete statement
// and returns it along with any error that occurred during the build process.
func (b *deleteBuilder) Build() (string, error) {
	return buildDeleteFrom(b)
}

// ToCQL builds the CQL query string for the delete statement.
func (b *deleteBuilder) ToCQL() string {
	cql, _ := buildDeleteFrom(b)
	return cql
}

func buildDeleteFrom(b *deleteBuilder) (string, error) {
	var sb strings.Builder

	sb.WriteString("DELETE")
	for i := range b.columns {
		if i > 0 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(" ")
		}
		sb.WriteString(b.columns[i])
	}
	sb.WriteString(" FROM ")
	sb.WriteString(b.table)

	// TODO(tjons): refactor this into a common function if possible
	for i := range b.using {
		if i == 0 {
			sb.WriteString(usingFragment)
		} else {
			sb.WriteString(andFragment)
		}
		sb.WriteString(string(b.using[i].T))
		sb.WriteString(spaceFragment)
		sb.WriteString(strconv.FormatInt(b.using[i].Arg, 10))
	}

	for i := range b.filterTerms {
		if i == 0 {
			sb.WriteString(" WHERE ")
		} else {
			sb.WriteString(" AND ")
		}

		sb.WriteString(b.filterTerms[i].column)
		sb.WriteString(" ")
		sb.WriteString(string(b.filterTerms[i].operator))
		sb.WriteString(" ")

		switch {
		case b.filterTerms[i].value != nil:
			sb.WriteString("?")
			b.queryValues = append(b.queryValues, b.filterTerms[i].value)
		case b.filterTerms[i].values != nil:
			sb.WriteString("(")
			for j := range b.filterTerms[i].values {
				if j > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString("?")
				b.queryValues = append(b.queryValues, b.filterTerms[i].values[j])
			}
			sb.WriteString(")")
		case b.filterTerms[i].deepValues != nil:
			sb.WriteString("(")
			for j := range b.filterTerms[i].deepValues {
				if j > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString("?")
				b.queryValues = append(b.queryValues, b.filterTerms[i].deepValues[j])
			}
			sb.WriteString(")")
		}
	}

	if b.ifExists {
		sb.WriteString(" IF EXISTS")
	}

	return sb.String(), nil
}

// QueryValues returns the values to be used in the query,
// in the order they were provided.
func (b *deleteBuilder) QueryValues() []any {
	return b.queryValues
}
