package service

import (
	"context"
	"time"

	"github.com/abyanjksatu/goscription/internal/database/mysql"
	"github.com/abyanjksatu/goscription/models"
	"github.com/abyanjksatu/goscription/util"
)

//UserService represent the service of the article
type UserService interface {
	// Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, csr string, err error)
	GetByID(ctx context.Context, id int64) (res models.User, err error)
	// Update(context.Context, ArticleParam) (err error)
	GetByEmail(ctx context.Context, email string) (res models.User, err error)
	Store(context.Context, UserParam) (res models.User, err error)
	// Delete(ctx context.Context, id int64) (err error)
}

//UserParam is paramter for Store Param
type UserParam struct {
	ID int64 `json:"id"`
	Name   string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type userService struct {
	userRepo    mysql.UserRepository
	contextTimeout time.Duration
}

// NewUserService will create new an articleService object representation of service.ArticleService interface
func NewUserService(a mysql.UserRepository, timeout time.Duration) UserService {
	if a == nil {
		panic("User repository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &userService{
		userRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *userService) Store(c context.Context, p UserParam) (res models.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedUser, _ := a.GetByEmail(ctx, p.Email)
	if existedUser != (models.User{}) {
		return models.User{}, util.ErrConflict
	}

	m := models.User{
		Name: p.Name,
		Email: p.Email,
		Password: p.Password,
	}

	err = a.userRepo.Store(ctx, &m)
	return
}

func (a *userService) GetByID(c context.Context, id int64) (res models.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.userRepo.GetByID(ctx, id)
	return
}

func (a *userService) GetByEmail(c context.Context, email string) (res models.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.userRepo.GetByEmail(ctx, email)
	return
}