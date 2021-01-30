package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/kecci/goscription/models"
	"github.com/kecci/goscription/util"
	"github.com/sirupsen/logrus"
)

// ArticleRepository represent the repository contract
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, csr string, err error)
	GetByID(ctx context.Context, id int64) (res models.Article, err error)
	GetByTitle(ctx context.Context, title string) (res models.Article, err error)
	Update(ctx context.Context, article *models.Article) (err error)
	Store(ctx context.Context, a *models.Article) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type mysqlArticleRepository struct {
	Conn *sql.DB
}

// NewArticleRepository will create an object that represent the article.Repository interface
func NewArticleRepository(conn *sql.DB) ArticleRepository {
	if conn == nil {
		panic("Database Connections is nil")
	}
	return &mysqlArticleRepository{conn}
}

func (m *mysqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []models.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result = make([]models.Article, 0)
	for rows.Next() {
		t := models.Article{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, nextCursor string, err error) {
	qbuilder := squirrel.Select("id", "title", "content", "updated_at", "created_at").From("article")
	qbuilder = qbuilder.OrderBy("id DESC").Limit(uint64(num))

	if cursor != "" {
		decodedCursor, err := strconv.ParseInt(cursor, 10, 64)
		if err != nil && cursor != "" {
			return nil, "", util.ErrBadParamInput
		}
		qbuilder = qbuilder.Where(squirrel.Lt{
			"id": decodedCursor,
		})
	}

	query, args, err := qbuilder.ToSql()
	if err != nil {
		return
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		return nil, "", err
	}

	nextCursor = cursor
	if len(res) > 0 {
		nextCursor = fmt.Sprintf("%d", res[len(res)-1].ID)
	}
	return
}

func (m *mysqlArticleRepository) GetByID(ctx context.Context, id int64) (res models.Article, err error) {
	query := `SELECT id,title,content, updated_at, created_at
  						FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return models.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, util.ErrNotFound
	}

	return
}

func (m *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (res models.Article, err error) {
	query := `SELECT id,title,content, updated_at, created_at
  						FROM article WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, util.ErrNotFound
	}
	return
}

func (m *mysqlArticleRepository) Store(ctx context.Context, a *models.Article) (err error) {
	query := `INSERT  article SET title=? , content=? , updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	res, err := stmt.ExecContext(ctx, a.Title, a.Content, now, now)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

func (m *mysqlArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		return
	}

	return
}
func (m *mysqlArticleRepository) Update(ctx context.Context, ar *models.Article) (err error) {
	query := `UPDATE article set title=?, content=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	now := time.Now()
	ar.UpdatedAt = now
	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, now, ar.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
		return
	}

	return
}
