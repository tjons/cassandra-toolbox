package qb

import (
	"strconv"
	"strings"
)

type InsertBuilder interface {
	QueryBuilder

	// Into specifies the table to insert into.
	Into(table string) InsertBuilder

	// Columns specifies the columns to insert values into.
	Columns(columns ...string) InsertBuilder

	// Values specifies the values to insert. The number of values provided
	// should match the number of columns specified and will be used in the
	// order they were provided.
	Values(values ...any) InsertBuilder

	// IfNotExists specifies that the INSERT query should include an IF NOT EXISTS clause.
	IfNotExists() InsertBuilder

	// Using specifies an optional USING clause for the INSERT query.
	// This can be used to specify a USING TIMESTAMP or USING TTL clause.
	// Multiple USING clauses can be specified, and they will be joined
	// with AND in the resulting query.
	// They will be added to the query in the order they were specified.
	Using(updateParam) InsertBuilder
}

type insertBuilder struct {
	table       string
	columns     []string
	using       []updateParam
	values      []any
	ifNotExists bool
}

// NewInsert creates a new InsertBuilder.
func NewInsert() InsertBuilder {
	return &insertBuilder{}
}

// Into specifies the table to insert into.
func (b *insertBuilder) Into(table string) InsertBuilder {
	b.table = table
	return b
}

// Using specifies an optional USING clause for the INSERT query.
// This can be used to specify a USING TIMESTAMP or USING TTL clause.
// Multiple USING clauses can be specified, and they will be joined
// with AND in the resulting query.
// They will be added to the query in the order they were specified.
func (b *insertBuilder) Using(param updateParam) InsertBuilder {
	b.using = append(b.using, param)
	return b
}

// Columns specifies the columns to insert values into.
func (b *insertBuilder) Columns(columns ...string) InsertBuilder {
	b.columns = append(b.columns, columns...)
	return b
}

// Values specifies the values to insert. The number of values provided
// should match the number of columns specified and will be used in the
// order they were provided.
func (b *insertBuilder) Values(values ...any) InsertBuilder {
	b.values = append(b.values, values...)
	return b
}

// IfNotExists specifies that the INSERT query should include an IF NOT EXISTS clause.
func (b *insertBuilder) IfNotExists() InsertBuilder {
	b.ifNotExists = true
	return b
}

// Build builds the CQL query string for the insert statement and
// returns it along with any error that occurred during the build process.
func (b *insertBuilder) Build() (string, error) {
	return buildInsertFrom(b)
}

// ToCQL builds the CQL query string for the insert statement.
func (b *insertBuilder) ToCQL() string {
	cql, _ := buildInsertFrom(b)
	return cql
}

// QueryValues returns the values to be used in the query,
// in the order they were provided.
func (b *insertBuilder) QueryValues() []any {
	vals := make([]any, 0, len(b.values))
	for _, v := range b.values {
		if _, ok := v.(literal); !ok {
			vals = append(vals, v)
		}
	}
	return vals
}

func buildInsertFrom(b *insertBuilder) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(b.table)
	sb.WriteString(" (")

	for i, col := range b.columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(col)
	}

	sb.WriteString(") VALUES (")

	for i := range b.values {
		if i > 0 {
			sb.WriteString(", ")
		}

		if literal, ok := b.values[i].(literal); ok {
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

	sb.WriteString(")")

	if b.ifNotExists {
		sb.WriteString(ifNotExistsFragment)
	}

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

	return sb.String(), nil
}
