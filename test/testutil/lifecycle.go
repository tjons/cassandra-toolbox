package testutil

import (
	"context"
	"fmt"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
)

func StartCassandraContainer(ctx context.Context, t *testing.T) (*cassandra.CassandraContainer, *gocql.Session, error) {
	cassandraContainer, err := cassandra.Run(ctx, "cassandra:5.0.6")
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

	if err := session.Query("SELECT release_version FROM system.local").Exec(); err != nil {
		t.Fatalf("failed to execute query: %v", err)
	}

	return cassandraContainer, session, nil
}
