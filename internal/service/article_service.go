package service

import (
	"context"
	"time"

	"github.com/abyanjksatu/goscription/internal/database/mysql"
	"github.com/abyanjksatu/goscription/models"
	"github.com/abyanjksatu/goscription/util"
)

//ArticleService represent the service of the article
type ArticleService interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, csr string, err error)
	GetByID(ctx context.Context, id int64) (res models.Article, err error)
	Update(ctx context.Context, ar *models.Article) (err error)
	GetByTitle(ctx context.Context, title string) (res models.Article, err error)
	Store(context.Context, *models.Article) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type articleService struct {
	articleRepo    mysql.ArticleRepository
	contextTimeout time.Duration
}

// NewArticleService will create new an articleService object representation of service.ArticleService interface
func NewArticleService(a mysql.ArticleRepository, timeout time.Duration) ArticleService {
	if a == nil {
		panic("Article repository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &articleService{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *articleService) Fetch(c context.Context, cursor string, num int64) (res []models.Article, nextCursor string, err error) {
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

func (a *articleService) GetByID(c context.Context, id int64) (res models.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	return
}

func (a *articleService) Update(c context.Context, ar *models.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleService) GetByTitle(c context.Context, title string) (res models.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	return
}

func (a *articleService) Store(c context.Context, m *models.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, m.Title)
	if existedArticle != (models.Article{}) {
		return util.ErrConflict
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *articleService) Delete(c context.Context, id int64) (err error) {
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
