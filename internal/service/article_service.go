package service

import (
	"context"
	"time"

	"github.com/kecci/goscription/internal/repository/mysql"
	"github.com/kecci/goscription/models"
	"github.com/kecci/goscription/util"
)

type (
	// ArticleService represent the service of the article
	ArticleService interface {
		Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, csr string, err error)
		GetByID(ctx context.Context, id int64) (res models.Article, err error)
		Update(context.Context, ArticleParam) (err error)
		GetByTitle(ctx context.Context, title string) (res models.Article, err error)
		Store(context.Context, ArticleParam) (err error)
		Delete(ctx context.Context, id int64) (err error)
	}

	// ArticleServiceImpl represent the service of the article
	ArticleServiceImpl struct {
		articleRepo    mysql.ArticleRepository
		contextTimeout time.Duration
	}
)

// ArticleParam ...
type ArticleParam struct {
	ID      int64  `json:"id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// NewArticleService will create new an articleService object representation of service.ArticleService interface
func NewArticleService(a mysql.ArticleRepository, timeout time.Duration) ArticleService {
	if a == nil {
		panic("Article repository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &ArticleServiceImpl{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

// Fetch ...
func (a *ArticleServiceImpl) Fetch(c context.Context, cursor string, num int64) (res []models.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

// GetByID ...
func (a *ArticleServiceImpl) GetByID(c context.Context, id int64) (res models.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	return
}

// Update ...
func (a *ArticleServiceImpl) Update(c context.Context, ap ArticleParam) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar := models.Article{
		ID:        ap.ID,
		Title:     ap.Title,
		Content:   ap.Content,
		UpdatedAt: time.Now(),
	}

	return a.articleRepo.Update(ctx, &ar)
}

// GetByTitle ...
func (a *ArticleServiceImpl) GetByTitle(c context.Context, title string) (res models.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	return
}

// Store ...
func (a *ArticleServiceImpl) Store(c context.Context, p ArticleParam) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, p.Title)
	if existedArticle != (models.Article{}) {
		return util.ErrConflict
	}

	m := models.Article{
		Title:   p.Title,
		Content: p.Content,
	}

	err = a.articleRepo.Store(ctx, &m)
	return
}

// Delete ...
func (a *ArticleServiceImpl) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (models.Article{}) {
		return util.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
