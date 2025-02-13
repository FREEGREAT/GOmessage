package database

import (
	"github.com/gocql/gocql"
)

type DBconnection struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

var connection DBconnection

func SetupDBConnection() (session *gocql.Session, err error) {
	connection.cluster = gocql.NewCluster("127.0.0.1:9042")
	connection.cluster.Consistency = gocql.Quorum
	connection.cluster.Keyspace = "chat_app"
	connection.session, err = connection.cluster.CreateSession()
	return connection.session, err
}
