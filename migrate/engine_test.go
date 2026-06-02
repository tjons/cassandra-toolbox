package migrate

import (
	"testing"
	"testing/fstest"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/tjons/cassandra-toolbox/test/testutil"
)

var sampleMigrations = fstest.MapFS{
	"001_create_table.cql": &fstest.MapFile{
		Data: []byte(
			`CREATE TABLE IF NOT EXISTS test.table1 (
				id uuid PRIMARY KEY,
				value text
			);`),
	},
}

func TestRunMigrations(t *testing.T) {
	t.Run("migrations run successfully", func(t *testing.T) {
		cass, session, err := testutil.StartCassandraContainer(t.Context(), t)
		if err != nil {
			t.Fatalf("failed to start Cassandra container: %v", err)
		}
		defer cass.Terminate(t.Context())
		defer session.Close()

		createTestKeyspace(t, session, "test")

		if err := RunMigrations(t.Context(), "test", session, sampleMigrations); err != nil {
			t.Fatalf("failed to run migrations: %v", err)
		}

		if err := session.Query("DESCRIBE TABLE test.table1").Exec(); err != nil {
			t.Fatalf("failed to verify migration: %v", err)
		}
	})
}

func createTestKeyspace(t *testing.T, session *gocql.Session, keyspace string) {
	if err := session.Query(
		`CREATE KEYSPACE IF NOT EXISTS ` + keyspace + ` WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}`,
	).Exec(); err != nil {
		t.Fatalf("failed to create keyspace: %v", err)
	}
}
