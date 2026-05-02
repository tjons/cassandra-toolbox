package qb_test

import (
	"testing"

	"github.com/tjons/cassandra-toolbox/qb"
)

func TestInsert(t *testing.T) {
	cases := []testCase{
		{
			name:           "Insert single column",
			expectedQuery:  "INSERT INTO test (column) VALUES (?)",
			builder:        qb.NewInsert().Into("test").Columns("column").Values("value"),
			expectedValues: []any{"value"},
		},
		{
			name:           "Insert multiple columns",
			expectedQuery:  "INSERT INTO test (column1, column2) VALUES (?, ?)",
			builder:        qb.NewInsert().Into("test").Columns("column1", "column2").Values("value1", "value2"),
			expectedValues: []any{"value1", "value2"},
		},
		{
			name:           "Conditional insert",
			expectedQuery:  "INSERT INTO test (column1, column2) VALUES (?, ?) IF NOT EXISTS",
			builder:        qb.NewInsert().Into("test").Columns("column1", "column2").Values("value1", "value2").IfNotExists(),
			expectedValues: []any{"value1", "value2"},
		},
		{
			name:           "Insert with raw CQL",
			expectedQuery:  "INSERT INTO test (column1, column2) VALUES (?, uuid())",
			builder:        qb.NewInsert().Into("test").Columns("column1", "column2").Values("value1", qb.CqlFunction("uuid()")),
			expectedValues: []any{"value1"},
		},
	}

	runTests(t, cases)
}
