package migratetest

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/tjons/cassandra-toolbox/migrate"
	"github.com/tjons/cassandra-toolbox/test/migratetest/schema"
)

func Migrate(ctx context.Context, keyspace string, session *gocql.Session) error {
	return migrate.RunMigrations(ctx, keyspace, session, schema.Files)
}
