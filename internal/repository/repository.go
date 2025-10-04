package repository

import (
	"FINAM_RADAR/internal/models"
	"context"
)

type INewsRepository interface {
	Add(news models.NewsRequest, ctx context.Context) (AddNews, error)
}
