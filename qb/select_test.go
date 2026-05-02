package qb_test

import (
	"testing"

	"github.com/tjons/cassandra-toolbox/qb"
)

func TestSelect(t *testing.T) {
	cases := []testCase{
		{
			name:           "Select all",
			expectedQuery:  `SELECT * FROM test`,
			builder:        qb.NewSelect().From("test"),
			expectedValues: []any{},
		},
		{
			name:           "Select single column",
			expectedQuery:  `SELECT entry_id FROM test`,
			builder:        qb.NewSelect().From("test").Column("entry_id"),
			expectedValues: []any{},
		},
		{
			name:           "Select with where clause",
			builder:        qb.NewSelect().From("test").Column("entry_id").Where("entry_id", qb.Equals("1")),
			expectedValues: []any{"1"},
			expectedQuery:  `SELECT entry_id FROM test WHERE entry_id = ?`,
		},
		{
			name:           "Select with where clause, IN operator, and multiple values",
			builder:        qb.NewSelect().From("test").Column("entry_id").Where("entry_id", qb.In("1", "2", "3")),
			expectedValues: []any{"1", "2", "3"},
			expectedQuery:  `SELECT entry_id FROM test WHERE entry_id IN (?, ?, ?)`,
		},
		{
			name:           "Select all where IN",
			builder:        qb.NewSelect().From("test").Where("column", qb.CollectionIn([]any{"a"}, []any{"a", "b"}, []any{"b"})),
			expectedValues: []any{[]any{"a"}, []any{"a", "b"}, []any{"b"}},
			expectedQuery:  `SELECT * FROM test WHERE column IN (?, ?, ?)`,
		},
	}

	runTests(t, cases)
}

func BenchmarkSelectAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queryStr, _ := qb.NewSelect().From("test").Build()
		if queryStr == "" {
			b.Fatal("Failed to build query")
		}
	}
}
