package qb_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tjons/cassandra-toolbox/qb"
)

type testCase struct {
	name           string
	expectedQuery  string
	builder        qb.QueryBuilder
	expectedValues []any
}

func runTests(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestUpdate: %s", tc.name), func(tt *testing.T) {
			qry, _ := tc.builder.Build()
			if qry == "" {
				tt.Fatal("Failed to build update query")
			}

			if qry != tc.expectedQuery {
				tt.Fatalf("Expected query %s, got %s", tc.expectedQuery, qry)
			}

			vals := tc.builder.QueryValues()
			if len(vals) == 0 && len(tc.expectedValues) > 0 {
				tt.Fatal("Failed to get query values")
			}

			if len(vals) != len(tc.expectedValues) {
				tt.Fatalf("Expected %d query values, got %d", len(tc.expectedValues), len(vals))
			}

			checkQueryValues(tt, tc.builder, tc.expectedValues...)
		})
	}
}

func checkQueryValues(t *testing.T, stmt qb.QueryBuilder, expectedVals ...any) {
	qvs := stmt.QueryValues()
	if len(qvs) != len(expectedVals) {
		t.Errorf("Expected %+v values in query, got %+v", qvs, expectedVals)
	}

	for i := range qvs {
		if !reflect.DeepEqual(qvs[i], expectedVals[i]) {
			t.Errorf(
				"Expected query elements to be equal, failed at index %d. %+v does not equal %+v",
				i, qvs[i], expectedVals[i])
		}
	}
}
