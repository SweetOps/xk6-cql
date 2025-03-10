package cql_test

import (
	"context"
	"testing"

	xk6_cql "github.com/SweetOps/xk6-cql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/cassandra"
)

var (
	CassandraImage = "cassandra:4.0.17"
)

func runCassandra(ctx context.Context, image string) (*cassandra.CassandraContainer, error) {
	return cassandra.Run(
		ctx,
		image,
	)
}

func Test_CQL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	c, err := runCassandra(ctx, CassandraImage)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Terminate(ctx) //nolint:errcheck // Error is not important here

	host, err := c.ConnectionHost(ctx)
	if err != nil {
		t.Fatal(err)
	}

	cql := xk6_cql.CQL{}

	err = cql.Session(
		xk6_cql.Config{
			Hosts:    []string{host},
			Keyspace: "system",
		},
	)
	require.NoError(t, err)

	err = cql.Exec("CREATE KEYSPACE IF NOT EXISTS test_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};") //nolint:lll // Long line
	require.NoError(t, err)

	err = cql.Session(xk6_cql.Config{
		Hosts:    []string{host},
		Keyspace: "test_keyspace",
	})
	require.NoError(t, err)

	err = cql.Exec("CREATE TABLE IF NOT EXISTS test_table (id INT PRIMARY KEY, name TEXT);")
	require.NoError(t, err)

	err = cql.Exec("INSERT INTO test_table (id, name) VALUES (1, 'test');")
	require.NoError(t, err)

	batchQueries := []string{
		"INSERT INTO test_table (id, name) VALUES (2, 'test2')",
		"INSERT INTO test_table (id, name) VALUES (3, 'test3')",
		"INSERT INTO test_table (id, name) VALUES (4, 'test4')",
	}

	err = cql.Batch("", batchQueries)
	assert.NoError(t, err)
}

func Test_CQL_Erros(t *testing.T) {
	t.Parallel()

	cql := xk6_cql.CQL{}
	err := cql.Session(xk6_cql.Config{})
	require.ErrorContains(t, err, "hosts and keyspace are required parameters")
}
