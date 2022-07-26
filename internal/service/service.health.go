package service

import "context"

type (
	// ArticleService represent the service of the article
	HealthService interface {
		CheckHealth(ctx context.Context) (err error)
	}

	// ArticleServiceImpl represent the service of the article
	HealthServiceImpl struct {
	}
)

func NewHealthService() HealthService {
	return &HealthServiceImpl{}
}

func (HealthServiceImpl) CheckHealth(ctx context.Context) (err error) {
	return nil
}
