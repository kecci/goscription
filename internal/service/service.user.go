package service

import (
	"context"
	"time"

	"github.com/kecci/goscription/internal/repository/mysql"
	"github.com/kecci/goscription/models"
	"github.com/kecci/goscription/utility"
)

type (
	// UserService represent the service of the article
	UserService interface {
		// Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, csr string, err error)
		GetByID(ctx context.Context, id int64) (res models.User, err error)
		// Update(context.Context, ArticleParam) (err error)
		GetByEmail(ctx context.Context, email string) (res models.User, err error)
		Store(context.Context, UserParam) (res models.User, err error)
		// Delete(ctx context.Context, id int64) (err error)
	}

	// UserServiceImpl represent the service of the article
	UserServiceImpl struct {
		userRepo       mysql.UserRepository
		contextTimeout time.Duration
	}
)

// UserParam ...
type UserParam struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// NewUserService will create new an articleService object representation of service.ArticleService interface
func NewUserService(a mysql.UserRepository, timeout time.Duration) UserService {
	return &UserServiceImpl{
		userRepo:       a,
		contextTimeout: timeout,
	}
}

// Store ...
func (a *UserServiceImpl) Store(c context.Context, p UserParam) (res models.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedUser, _ := a.GetByEmail(ctx, p.Email)
	if existedUser != (models.User{}) {
		return models.User{}, utility.ErrConflict
	}

	m := models.User{
		Name:     p.Name,
		Email:    p.Email,
		Password: p.Password,
	}

	err = a.userRepo.Store(ctx, &m)
	return
}

// GetByID ...
func (a *UserServiceImpl) GetByID(c context.Context, id int64) (res models.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.userRepo.GetByID(ctx, id)
	return
}

// GetByEmail ...
func (a *UserServiceImpl) GetByEmail(c context.Context, email string) (res models.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.userRepo.GetByEmail(ctx, email)
	return
}
