package qb_test

import (
	"fmt"
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
			if vals == nil {
				tt.Fatal("Failed to get query values")
			}

			if len(vals) != len(tc.expectedValues) {
				tt.Fatalf("Expected %d query values, got %d", len(tc.expectedValues), len(vals))
			}

			for i, v := range vals {
				if v != tc.expectedValues[i] {
					tt.Fatalf("Expected value %v, got %v", tc.expectedValues[i], v)
				}
			}
		})
	}
}
