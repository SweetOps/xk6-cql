package cql_test

import (
	"context"
	"testing"
	"time"

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
	time.Sleep(30 * time.Second)

	cql := xk6_cql.CQL{}

	err = cql.Session(xk6_cql.Config{
		Hosts:    []string{host},
		Keyspace: "system",
	})
	require.NoError(t, err)

	err = cql.Exec("SELECT * FROM local")
	assert.NoError(t, err)
}
