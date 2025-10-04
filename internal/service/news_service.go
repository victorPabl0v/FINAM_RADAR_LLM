package service

import (
	"FINAM_RADAR/internal/models"
	"FINAM_RADAR/internal/repository"
	"context"
	"strings"
)

type NewsService struct {
	INewsService
	newsRepo repository.INewsRepository
}

func NewNewsService(newsRepo repository.INewsRepository) *NewsService {
	return &NewsService{newsRepo: newsRepo}
}

func (this *NewsService) AddNews(news models.NewsRequest, ctx context.Context) (AddNewsResponse, error) {
	// normalize words to lower register
	for i, n := range news.News {
		for j, e := range n.Entities {
			news.News[i].Entities[j] = strings.ToLower(e)
		}
		for j, b := range n.Draft.Bullets {
			news.News[i].Draft.Bullets[j] = strings.ToLower(b)
		}
	}

	status, err := this.newsRepo.Add(news, ctx)

	var response AddNewsResponse
	if err != nil {
		response.Status = false
		response.Message = err.Error()
		return response, err
	}

	response.Status = status.Status
	if status.Status {
		response.Message = "success"
	} else {
		response.Message = "fail"
	}

	return response, nil
}
