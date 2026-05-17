package qb

import (
	"strconv"
	"strings"
)

// UpdateBuilder builds UPDATE queries.
type UpdateBuilder interface {
	QueryBuilder

	// Using specifies an optional USING clause for the UPDATE query.
	// This can be used to specify a USING TIMESTAMP or USING TTL clause.
	// Multiple USING clauses can be specified, and they will be joined with AND in the resulting query.
	// They will be added to the query in the order they were specified.
	Using(updateParam) UpdateBuilder

	// Table specifies the table to update.
	Table(name string) UpdateBuilder

	// Set specifies a column to update and the value to update it to. Multiple calls to Set will add multiple columns to the SET clause of the UPDATE query.
	Set(column string, value any) UpdateBuilder

	// Where specifies a condition for the UPDATE query. Multiple calls to Where will be joined with AND in the resulting query.
	Where(condition string, value filterTerm) UpdateBuilder

	// IfExists specifies that the UPDATE query should include an IF EXISTS clause.
	IfExists() UpdateBuilder
}

type updateBuilder struct {
	table        string
	using        []updateParam
	columns      []string
	columnValues []any
	conditions   []filterTerm
	ifExists     bool
	queryValues  []any
}

// NewUpdate creates a new UpdateBuilder.
func NewUpdate() UpdateBuilder {
	return &updateBuilder{}
}

// Table specifies the table to update.
func (b *updateBuilder) Table(name string) UpdateBuilder {
	b.table = name
	return b
}

// Using specifies an optional USING clause for the UPDATE query. This can be used to specify a USING TIMESTAMP or USING TTL clause.
// Multiple USING clauses can be specified, and they will be joined with AND in the resulting query.
// They will be added to the query in the order they were specified.
func (b *updateBuilder) Using(param updateParam) UpdateBuilder {
	b.using = append(b.using, param)
	return b
}

// Set specifies a column to update and the value to update it to. Multiple calls to Set will add multiple columns to the SET clause of the UPDATE query.
func (b *updateBuilder) Set(column string, value any) UpdateBuilder { // TODO(tjons): this needs to be able to accept more than just a single value, like setting a collection type...
	b.columns = append(b.columns, column)
	b.columnValues = append(b.columnValues, value)
	return b
}

// Where specifies a condition for the UPDATE query. Multiple calls to Where will be joined with AND in the resulting query.
func (b *updateBuilder) Where(column string, ft filterTerm) UpdateBuilder {
	ft.column = column
	b.conditions = append(b.conditions, ft)

	return b
}

func (b *updateBuilder) IfExists() UpdateBuilder {
	b.ifExists = true
	return b
}

// Build builds the UPDATE query and returns the query string and any error that occurred during building.
func (b *updateBuilder) Build() (string, error) {
	return buildUpdateFrom(b)
}

// ToCQL builds the UPDATE query and returns the query string without
// any error validation. If you want error validation, use Build instead.
func (b *updateBuilder) ToCQL() string {
	cql, _ := buildUpdateFrom(b)
	return cql
}

// QueryValues returns the values to be used in the query, in the order they were provided.
func (b *updateBuilder) QueryValues() []any {
	// TODO(tjons): clean this up, I don't think skippedLiteral* is as relevant for the UPDATE clause

	// preallocate slice with the most elements we might need
	vals := make([]any, len(b.columnValues)+len(b.conditions))
	if len(vals) == 0 {
		return nil
	}

	var skippedLiteralColumns, skippedLiteralConditions int
	for i := range b.columnValues {
		if _, ok := b.columnValues[i].(literal); ok {
			skippedLiteralColumns++
		} else {
			vals[i-skippedLiteralColumns] = b.columnValues[i]
		}
	}

	offset := len(b.columnValues) - skippedLiteralColumns
	for i := range b.conditions {
		vals[offset+i-skippedLiteralConditions] = b.conditions[i].value
	}

	// chop off unused elements
	return vals[:len(vals)-skippedLiteralColumns-skippedLiteralConditions]
}

func buildUpdateFrom(b *updateBuilder) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("UPDATE ")
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

	sb.WriteString(" SET ")
	for i := range b.columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(b.columns[i])
		sb.WriteString(" = ")
		if literal, ok := b.columnValues[i].(literal); ok {
			if literal.singleQuote {
				sb.WriteString("'")
			}
			sb.WriteString(literal.value)
			if literal.singleQuote {
				sb.WriteString("'")
			}
		} else {
			sb.WriteString("?")
		}
	}

	sb.WriteString(" WHERE ")

	for i := range b.conditions {
		if i > 0 {
			sb.WriteString(" AND ")
		}

		sb.WriteString(b.conditions[i].column)
		sb.WriteString(" = ")

		if len(b.conditions[i].values) > 0 {
			sb.WriteString("(")
			for j := range b.conditions[i].values {
				if j > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString("?")
			}
			sb.WriteString(")")
		} else {
			sb.WriteString("?")
		}
	}

	if b.ifExists {
		// TODO(tjons): support conditional check
		sb.WriteString(" IF EXISTS")
	}

	return sb.String(), nil
}
