package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kecci/goscription/internal/service"
	"github.com/kecci/goscription/mocks"
	"github.com/kecci/goscription/models"
	"github.com/kecci/goscription/util"
)

func TestFetch(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := models.Article{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]models.Article, 0)
	mockListArtilce = append(mockListArtilce, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		u := service.NewArticleService(mockArticleRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockArticleRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)

	})

}

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := models.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)

	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Article{}, errors.New("Unexpected")).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, models.Article{}, a)

		mockArticleRepo.AssertExpectations(t)

	})

}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticleParam := service.ArticleParam{
		Title:   "Hello",
		Content: "Content",
	}

	mockArticle := models.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticleParam

		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(models.Article{}, util.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticleParam.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		err := u.Store(context.TODO(), mockArticleParam)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)

	})

}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := models.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()

		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)

	})
	t.Run("article-is-not-exist", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Article{}, nil).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)

	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(models.Article{}, errors.New("Unexpected Error")).Once()

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)

	})

}

func TestUpdate(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticleParam := service.ArticleParam{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Article")).Once().Return(nil)

		u := service.NewArticleService(mockArticleRepo, time.Second*2)

		err := u.Update(context.TODO(), mockArticleParam)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
