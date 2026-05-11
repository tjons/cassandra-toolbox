package qb_test

import (
	"testing"

	"github.com/tjons/cassandra-toolbox/qb"
)

func TestUpdate(t *testing.T) {
	testCases := []testCase{
		{
			name:           "Update single column",
			expectedQuery:  "UPDATE test SET column = ? WHERE id = ?",
			builder:        qb.NewUpdate().Table("test").Set("column", "value").Where("id", qb.Equals("1")),
			expectedValues: []any{"value", "1"},
		},
		{
			name:           "Update multiple columns, single condition",
			expectedQuery:  "UPDATE test SET column1 = ?, column2 = ? WHERE id = ?",
			builder:        qb.NewUpdate().Table("test").Set("column1", "value1").Set("column2", "value2").Where("id", qb.Equals("1")),
			expectedValues: []any{"value1", "value2", "1"},
		},
		{
			name:           "Update multiple columns, multiple conditions",
			expectedQuery:  "UPDATE test SET column1 = ?, column2 = ? WHERE id = ? AND status = ?",
			builder:        qb.NewUpdate().Table("test").Set("column1", "value1").Set("column2", "value2").Where("id", qb.Equals("1")).Where("status", qb.Equals("active")),
			expectedValues: []any{"value1", "value2", "1", "active"},
		},
		{
			name:           "Update with if exists",
			expectedQuery:  "UPDATE test SET column = ? WHERE id = ? IF EXISTS",
			builder:        qb.NewUpdate().Table("test").Set("column", "value").Where("id", qb.Equals("1")).IfExists(),
			expectedValues: []any{"value", "1"},
		},
		{
			name:           "Update with USING TTL",
			expectedQuery:  "UPDATE test USING TTL 3600 SET column = ? WHERE id = ?",
			builder:        qb.NewUpdate().Table("test").Using(qb.TTL(3600)).Set("column", 1).Where("id", qb.Equals(1)),
			expectedValues: []any{1, 1},
		},
	}

	runTests(t, testCases)
}
