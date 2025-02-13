package pkg

import (
	"crypto/tls"
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func ConnectDB() *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{"qywfvq4856.us-east1.gcp.clickhouse.cloud:9440"}, // 9440 is a secure native TCP port
		Protocol: clickhouse.Native,
		TLS:      &tls.Config{}, // enable secure TLS
		Auth: clickhouse.Auth{
			Username: "default",
			Password: "mzEonKSd_9X7B",
		},
	})

	return conn
}
