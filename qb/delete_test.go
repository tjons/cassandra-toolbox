package qb_test

import (
	"testing"

	"github.com/tjons/cassandra-toolbox/qb"
)

func TestDelete(t *testing.T) {
	cases := []testCase{
		{
			name:           "Delete with single condition",
			expectedQuery:  "DELETE FROM test WHERE id = ?",
			builder:        qb.NewDelete().From("test").Where("id", qb.Equals("1")),
			expectedValues: []any{"1"},
		},
		{
			name:           "Delete with multiple conditions",
			expectedQuery:  "DELETE FROM test WHERE id = ? AND name = ?",
			builder:        qb.NewDelete().From("test").Where("id", qb.Equals("1")).Where("name", qb.Equals("Alice")),
			expectedValues: []any{"1", "Alice"},
		},
		{
			name:           "Delete with IF EXISTS",
			expectedQuery:  "DELETE FROM test WHERE id = ? IF EXISTS",
			builder:        qb.NewDelete().From("test").Where("id", qb.Equals("1")).IfExists(),
			expectedValues: []any{"1"},
		},
		{
			name:           "Delete with single column",
			expectedQuery:  "DELETE column FROM test WHERE id = ?",
			builder:        qb.NewDelete().Column("column").From("test").Where("id", qb.Equals("1")),
			expectedValues: []any{"1"},
		},
		{
			name:           "Delete with multiple columns",
			expectedQuery:  "DELETE column1, column2 FROM test WHERE id = ?",
			builder:        qb.NewDelete().Column("column1").Column("column2").From("test").Where("id", qb.Equals("1")),
			expectedValues: []any{"1"},
		},
	}

	runTests(t, cases)
}
