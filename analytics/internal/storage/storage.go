package storage

import (
	"context"

	"gomessage.com/analytics/internal/models"
)

type AnalyticsRepository interface {
	AddData(ctx context.Context, user *models.Analytics) (string, error)
}
