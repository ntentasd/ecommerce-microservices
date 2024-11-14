package database

import (
	"fmt"

	"github.com/gocql/gocql"
)

type CassandraDB struct {
	session *gocql.Session
}

func NewCassandraDB() (*CassandraDB, error) {
	cluster := gocql.NewCluster("127.0.0.1") // Update with your Cassandra host
	cluster.Keyspace = "ecommerce"           // Your keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cassandra: %w", err)
	}
	return &CassandraDB{session: session}, nil
}
