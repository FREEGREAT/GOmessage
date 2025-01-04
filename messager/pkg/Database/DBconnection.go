package database

import (
	"github.com/gocql/gocql"
)

type DBconnection struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

var connection DBconnection

func SetupDBConnection() (err error) {
	connection.cluster = gocql.NewCluster("127.0.0.1:9042")
	connection.cluster.Consistency = gocql.Quorum
	connection.cluster.Keyspace = "messager"
	connection.session, err = connection.cluster.CreateSession()
	return err
}
