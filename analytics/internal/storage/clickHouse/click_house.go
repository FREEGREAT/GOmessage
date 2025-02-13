package clickhouse

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	"gomessage.com/analytics/internal/models"
	"gomessage.com/analytics/internal/storage"
	"gomessage.com/analytics/pkg"
)

const (
	transactionSuccesfuly = "OK"
	transactionERROR      = "ERROR"
	testDataRegion        = "UA"
	testDataSource        = "web"
)

type analyticsRepository struct {
	conn *sql.DB
}

func (a *analyticsRepository) AddData(ctx context.Context, user *models.Analytics) (string, error) {
	if err := pkg.InitConfig(); err != nil {
		panic("Error while initialising config!(analytics)")
	}
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	row := a.conn.QueryRow("INSERT INTO user_registration_metrics(date, datetime, user_id, registration_source, country) VALUES (?,?,?,?,?)", currentDate, time.Now(), user.ID, user.Source, user.Region)
	if row.Err() != nil {
		str := row.Err()
		logrus.Errorf("Error while updating analytic DB: %d", str)
		return transactionERROR, row.Err()
	}
	return transactionSuccesfuly, nil
}

func NewAnalyticsRepository(connection *sql.DB) storage.AnalyticsRepository {
	return &analyticsRepository{
		conn: connection,
	}
}
