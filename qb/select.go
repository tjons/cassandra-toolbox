package qb

import (
	"strconv"
	"strings"
)

// SelectBuilder builds SELECT queries.
type SelectBuilder interface {
	QueryBuilder

	// Distinct specifies that the SELECT query should include a DISTINCT clause.
	Distinct() SelectBuilder

	// Json specifies that the SELECT query should include a JSON clause, which will return results as JSON objects instead of rows.
	Json() SelectBuilder

	// ColumnAs specifies a column to select and an optional alias for that column.
	// If an alias is provided, the resulting query will include an AS clause for that column.
	ColumnAs(string, string) SelectBuilder

	// Column specifies a column to select. Multiple calls to Column will add multiple columns to the SELECT clause.
	Column(string) SelectBuilder

	// Columns specifies multiple columns to select. It is equivalent to calling Column for each column.
	Columns([]string) SelectBuilder

	// From specifies the table to select from.
	From(string) SelectBuilder

	// Where specifies a condition for the SELECT query. Multiple calls to Where will be joined with AND in the resulting query.
	Where(string, filterTerm) SelectBuilder

	// PerPartitionLimit specifies a PER PARTITION LIMIT for the SELECT query.
	PerPartitionLimit(uint) SelectBuilder

	// Limit specifies a LIMIT for the SELECT query.
	Limit(uint) SelectBuilder

	// OrderBy specifies one or more columns to order the results by, along with the direction (ASC or DESC) for each column.
	OrderBy(...orderByClause) SelectBuilder

	// GroupBy specifies one or more columns to group the results by.
	GroupBy(...string) SelectBuilder

	// AllowFiltering specifies that the SELECT query should include an ALLOW FILTERING clause.
	AllowFiltering() SelectBuilder
}

// NewSelect creates a new SelectBuilder.
func NewSelect() SelectBuilder {
	return &selectBuilder{}
}

type selectBuilder struct {
	cols []string

	// retrieveColumns is a map of requested columns: potential aliases.
	// if no alias is set for a given column, the key will have an empty
	// value.
	retrieveColumns   map[string]string
	table             string
	filterTerms       []*filterTerm
	limit             uint
	perPartitionLimit uint
	values            []any
	queryValues       []any
	allowFiltering    bool
	isDistinct        bool
	isJson            bool
	orderBy           []orderByClause
	groupBy           []string
}

// Json specifies that the SELECT query should include a JSON clause, which will return results as JSON objects instead of rows.
func (b *selectBuilder) Json() SelectBuilder {
	b.isJson = true

	return b
}

// PerPartitionLimit specifies a PER PARTITION LIMIT for the SELECT query.
func (b *selectBuilder) PerPartitionLimit(num uint) SelectBuilder {
	b.perPartitionLimit = num

	return b
}

// Column specifies a column to select. Multiple calls to Column will add multiple columns to the SELECT clause.
func (b *selectBuilder) Column(name string) SelectBuilder {
	if b.retrieveColumns == nil {
		b.retrieveColumns = make(map[string]string)
	}

	if _, exists := b.retrieveColumns[name]; exists {
		return b
	}

	b.retrieveColumns[name] = ""
	b.cols = append(b.cols, name)

	return b
}

// ColumnAs specifies a column to select and an optional alias for that column.
func (b *selectBuilder) ColumnAs(name, alias string) SelectBuilder {
	if b.retrieveColumns == nil {
		b.retrieveColumns = make(map[string]string)
	}

	if setAlias, exists := b.retrieveColumns[name]; exists && setAlias == alias {
		return b
	}

	b.retrieveColumns[name] = alias
	b.cols = append(b.cols, name)

	return b
}

// Columns specifies multiple columns to select. It is equivalent to calling Column for each column.
func (b *selectBuilder) Columns(names []string) SelectBuilder {
	if b.retrieveColumns == nil {
		b.retrieveColumns = make(map[string]string)
	}

	for _, name := range names {
		b.Column(name)
	}

	return b
}

// From specifies the table to select from.
func (b *selectBuilder) From(table string) SelectBuilder {
	b.table = table

	return b
}

// Where specifies a condition for the SELECT query. Multiple calls to Where will be joined with AND in the resulting query.
func (b *selectBuilder) Where(column string, ft filterTerm) SelectBuilder {
	ft.column = column

	b.filterTerms = append(b.filterTerms, &ft)

	return b
}

// Limit specifies a LIMIT for the SELECT query.
func (b *selectBuilder) Limit(num uint) SelectBuilder {
	b.limit = num

	return b
}

// Build builds the CQL query string for the select statement and returns it along with any error that occurred during the build process.
func (b *selectBuilder) Build() (string, error) {
	return buildSelectFrom(b)
}

// ToCQL builds the CQL query string for the select statement.
func (b *selectBuilder) ToCQL() string {
	cql, _ := buildSelectFrom(b)
	return cql
}

func buildSelectFrom(b *selectBuilder) (string, error) {
	q := strings.Builder{}
	q.WriteString(selectFragment)
	if b.isDistinct {
		q.WriteString(distinctFragment)
	} else if b.isJson {
		q.WriteString(jsonFragment)
	}

	for i, col := range b.cols { // ordinal is important for scanning
		if i > 0 {
			q.WriteString(", ")
		}
		q.WriteString(col)
		if alias, ok := b.retrieveColumns[col]; ok && alias != "" {
			q.WriteString(" AS ")
			q.WriteString(alias)
		}
	}
	if len(b.cols) == 0 {
		q.WriteString("*")
	}

	q.WriteString(fromFragment)
	q.WriteString(b.table)

	for i := range b.filterTerms {
		if i == 0 {
			q.WriteString(whereFragment)
		} else {
			q.WriteString(andFragment)
		}

		q.WriteString(b.filterTerms[i].column)
		q.WriteString(spaceFragment)
		q.WriteString(string(b.filterTerms[i].operator))
		q.WriteString(spaceFragment)

		switch {
		case b.filterTerms[i].value != nil:
			q.WriteString("?")
			b.queryValues = append(b.queryValues, b.filterTerms[i].value)
		case b.filterTerms[i].values != nil:
			q.WriteString("(")

			for j := range b.filterTerms[i].values {
				if j > 0 {
					q.WriteString(", ")
				}
				q.WriteString("?")
				b.queryValues = append(b.queryValues, b.filterTerms[i].values[j])
			}
			q.WriteString(")")
		case b.filterTerms[i].deepValues != nil:
			q.WriteString("(")
			for j := range b.filterTerms[i].deepValues {
				if j > 0 {
					q.WriteString(", ")
				}
				q.WriteString("?")
				b.queryValues = append(b.queryValues, b.filterTerms[i].deepValues[j])
			}
			q.WriteString(")")
		}
	}

	if len(b.groupBy) > 0 {
		q.WriteString(groupByFragment)

		for i, col := range b.groupBy {
			if i > 0 {
				q.WriteString(", ")
			}
			q.WriteString(col)
		}
	}

	if len(b.orderBy) > 0 {
		q.WriteString(orderByFragment)

		for i, term := range b.orderBy {
			if i > 0 {
				q.WriteString(", ")
			}
			q.WriteString(term.column)
			q.WriteString(" ")
			q.WriteString(string(term.dir))
		}
	}

	if b.perPartitionLimit > 0 {
		q.WriteString(perPartitionLimitFragment)
		q.WriteString(strconv.Itoa(int(b.perPartitionLimit)))
	}

	if b.limit > 0 {
		q.WriteString(limitFragment)
		q.WriteString(strconv.Itoa(int(b.limit)))
	}

	if b.allowFiltering {
		q.WriteString(allowFilteringFragment)
	}

	// TODO(tjons): grab these builders from a mempool
	// also, validate
	return q.String(), nil
}

// AllowFiltering specifies that the SELECT query should include an ALLOW FILTERING clause.
func (b *selectBuilder) AllowFiltering() SelectBuilder {
	b.allowFiltering = true
	return b
}

// Distinct specifies that the SELECT query should include a DISTINCT clause.
func (b *selectBuilder) Distinct() SelectBuilder {
	b.isDistinct = true
	return b
}

// QueryValues returns the values to be used in the query, in the order they were provided.
func (b *selectBuilder) QueryValues() []any {
	return b.queryValues
}

// GroupBy specifies one or more columns to group the results by.
func (b *selectBuilder) GroupBy(cols ...string) SelectBuilder {
	b.groupBy = cols

	return b
}

// OrderBy specifies one or more columns to order the results by, along with the direction (ASC or DESC) for each column.
func (b *selectBuilder) OrderBy(terms ...orderByClause) SelectBuilder {
	b.orderBy = terms

	return b
}

type direction string

const (
	// Asc specifies that the results should be ordered in ascending order.
	Asc direction = "ASC"

	// Desc specifies that the results should be ordered in descending order.
	Desc direction = "DESC"
)

type orderByClause struct {
	column string
	dir    direction
}

// Writetime applies the WRITETIME function to a column, which returns the timestamp of when the column was last written to.
func Writetime(column string) string {
	sb := strings.Builder{}
	sb.WriteString(writetimeFragment)
	sb.WriteString(openParenFragment)
	sb.WriteString(column)
	sb.WriteString(closeParenFragment)

	return sb.String()
}

// MaxWritetime applies the MAXWRITETIME function to a column, which returns the maximum timestamp of when any cell in the column was last written to.
func MaxWritetime(column string) string {
	sb := strings.Builder{}
	sb.WriteString(maxWritetimeFragment)
	sb.WriteString(openParenFragment)
	sb.WriteString(column)
	sb.WriteString(closeParenFragment)

	return sb.String()
}

// Ttl applies the TTL function to a column, which returns the remaining time to live for the row value in seconds.
func Ttl(column string) string {
	sb := strings.Builder{}
	sb.WriteString(ttlFragment)
	sb.WriteString(openParenFragment)
	sb.WriteString(column)
	sb.WriteString(closeParenFragment)

	return sb.String()
}
