package service

import (
	"FINAM_RADAR/internal/models"
	"context"
)

type INewsService interface {
	AddNews(news models.NewsRequest, ctx context.Context) (AddNewsResponse, error)
}
