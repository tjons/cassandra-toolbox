package test

import (
	"fmt"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
	"github.com/tjons/cassandra-toolbox/qb"
	"github.com/tjons/cassandra-toolbox/test/migratetest"
)

func TestCassandra(t *testing.T) {
	ctx := t.Context()
	cassandraContainer, err := cassandra.Run(ctx, "cassandra:5.0.6")
	defer func() {
		if err := cassandraContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate Cassandra container: %v", err)
		}
	}()
	if err != nil {
		t.Fatalf("failed to start Cassandra container: %v", err)
	}

	port, err := cassandraContainer.MappedPort(ctx, "9042")
	if err != nil {
		t.Fatalf("failed to get Cassandra container port: %v", err)
	}

	cluster := gocql.NewCluster(fmt.Sprintf("localhost:%s", port.Port()))
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		t.Fatalf("failed to create Cassandra session: %v", err)
	}
	defer session.Close()

	if err := session.Query("SELECT release_version FROM system.local").Exec(); err != nil {
		t.Fatalf("failed to execute query: %v", err)
	}

	t.Run("run migrations", func(t *testing.T) {
		if err := session.Query("CREATE KEYSPACE IF NOT EXISTS test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}").Exec(); err != nil {
			t.Fatalf("failed to create keyspace: %v", err)
		}

		if err := migratetest.Migrate(ctx, "test", session); err != nil {
			t.Fatalf("failed to migrate keyspace: %v", err)
		}
	})

	t.Run("select with IN filter on secondary index", func(t *testing.T) {
		insertQuery := qb.NewInsert().
			Into("test.table1").
			Columns("entity_name", "entity_attribute_1").
			Values("test_entity", "test_attribute")
		iq, _ := insertQuery.Build()

		if err := session.Query(iq, insertQuery.QueryValues()...).Exec(); err != nil {
			t.Fatalf("failed to insert data: %v", err)
		}

		set := []string{"test_attribute", "non_existent_attribute"}
		selectQuery := qb.NewSelect().
			Column("entity_name").
			From("test.table1").
			Where("entity_attribute_1", qb.In(set...)).
			AllowFiltering()
		sq, _ := selectQuery.Build()

		var results []string
		scanner := session.Query(sq, selectQuery.QueryValues()...).Iter().Scanner()
		for scanner.Next() {
			var name string
			if err := scanner.Scan(&name); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			results = append(results, name)
		}
		if err := scanner.Err(); err != nil {
			t.Fatalf("failed to iterate over results: %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
	})
}
