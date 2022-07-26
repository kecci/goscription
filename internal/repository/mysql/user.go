package mysql

import (
	"context"
	"database/sql"

	"github.com/kecci/goscription/internal/library/db"
	"github.com/kecci/goscription/models"
	"github.com/kecci/goscription/utility"
	"github.com/sirupsen/logrus"
)

// UserRepository represent the repository contract
type UserRepository interface {
	Store(ctx context.Context, a *models.User) (err error)
	// Get(ctx context.Context, cursor string, num int64) (res []models.Article, csr string, err error)
	GetByID(ctx context.Context, id int64) (res models.User, err error)
	GetByEmail(ctx context.Context, email string) (res models.User, err error)
	// Update(ctx context.Context, article *models.Article) (err error)
	// Delete(ctx context.Context, id int64) (err error)
}

type mysqlUserRepository struct {
	Conn *sql.DB
}

// NewUserRepository will create an object that represent the article.Repository interface
func NewUserRepository(DB db.Database) UserRepository {
	if DB.Mysql == nil {
		panic("Database Connections is nil")
	}
	return &mysqlUserRepository{DB.Mysql}
}

func (m *mysqlUserRepository) Store(ctx context.Context, a *models.User) (err error) {
	query := `INSERT user SET name=?, email=?, password=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Name, a.Email, a.Password)
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

func (m *mysqlUserRepository) GetByID(ctx context.Context, id int64) (res models.User, err error) {
	query := `SELECT * FROM user WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return models.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, utility.ErrNotFound
	}

	return
}

func (m *mysqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []models.User, err error) {
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

	result = make([]models.User, 0)
	for rows.Next() {
		t := models.User{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.Password,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, title string) (res models.User, err error) {
	query := `SELECT * FROM user WHERE email = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, utility.ErrNotFound
	}
	return
}
