package mysql_test

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"testing"
	"time"

	repoHandler "github.com/kecci/goscription/internal/repository/mysql"
	"github.com/kecci/goscription/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mysqlArticleSuiteTest struct {
	MysqlSuite
}

func TestArticleSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip article mysql repository test")
	}
	dsn := os.Getenv("MYSQL_TEST_URL")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:33060)/testing?parseTime=1&loc=Asia%2FJakarta&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	}
	articleSuite := &mysqlArticleSuiteTest{
		MysqlSuite{
			DSN:                     dsn,
			MigrationLocationFolder: "migrations",
		},
	}

	suite.Run(t, articleSuite)
}

func (s *mysqlArticleSuiteTest) SetupTest() {
	log.Println("Starting a Test. Migrating the Database")
	err, _ := s.Migration.Up()
	require.NoError(s.T(), err)
	log.Println("Database Migrated Successfully")
}

func (s *mysqlArticleSuiteTest) TearDownTest() {
	log.Println("Finishing Test. Dropping The Database")
	err, _ := s.Migration.Down()
	require.NoError(s.T(), err)
	log.Println("Database Dropped Successfully")
}

// https://blevesearch.com/news/Deferred-Cleanup,-Checking-Errors,-and-Potential-Problems/
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
func getArticleByID(t *testing.T, DBconn *sql.DB, id int64) models.Article {
	var res models.Article

	query := `SELECT id, title, content, created_at, updated_at FROM article WHERE id=?`

	row := DBconn.QueryRow(query, id)
	err := row.Scan(
		&res.ID,
		&res.Title,
		&res.Content,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err == nil {
		return res
	}

	if err != sql.ErrNoRows {
		require.NoError(t, err)
	}

	return res
}

func getMockArrArticle() []models.Article {
	return []models.Article{
		models.Article{
			ID:      1,
			Title:   "Tekno",
			Content: "tekno",
		},
		models.Article{
			ID:      2,
			Title:   "Bola",
			Content: "bola",
		},
		models.Article{
			ID:      3,
			Title:   "Asmara",
			Content: "asmara",
		},
		models.Article{
			ID:      4,
			Title:   "Celebs",
			Content: "celebs",
		},
	}
}

func seedArticleData(t *testing.T, DBConn *sql.DB) {
	arrArticles := getMockArrArticle()
	query := `INSERT article SET id=?, title=?, content=?, created_at=?, updated_at=?`
	stmt, err := DBConn.Prepare(query)
	require.NoError(t, err)
	defer Close(stmt)
	for _, article := range arrArticles {
		_, err := stmt.Exec(article.ID, article.Title, article.Content, time.Now(), time.Now())
		require.NoError(t, err)
	}
}

func (m *mysqlArticleSuiteTest) TestStore() {

	// Prepare Steps
	repo := repoHandler.NewArticleRepository(m.DBConn)

	type testCase struct {
		Title          string
		Payload        *models.Article
		ExpectedResult error
	}
	now := time.Now()
	arrTestcase := []testCase{
		testCase{
			Title: "store-success",
			Payload: &models.Article{
				Title:     "News",
				Content:   "news",
				CreatedAt: now,
				UpdatedAt: now,
			},
			ExpectedResult: nil,
		},
	}

	for _, tc := range arrTestcase {
		m.T().Run(tc.Title, func(t *testing.T) {
			err := repo.Store(context.Background(), tc.Payload)
			require.Equal(m.T(), tc.ExpectedResult, err)
			if err == nil {
				assert.NotZero(m.T(), tc.Payload.ID)
				res := getArticleByID(m.T(), m.DBConn, tc.Payload.ID)
				assert.NotNil(m.T(), res)
				assert.Equal(m.T(), tc.Payload.Content, res.Content)
			}
		})
	}
}

func (m *mysqlArticleSuiteTest) TestFetch() {
	repo := repoHandler.NewArticleRepository(m.DBConn)
	seedArticleData(m.T(), m.DBConn)

	type testCase struct {
		Title          string
		Num            int64
		Cursor         string
		ExpectedResult []models.Article
	}

	arrTestcase := []testCase{
		testCase{
			Title: "fetch-without-cursor-success",
			Num:   3,
			ExpectedResult: []models.Article{
				models.Article{
					ID:      4,
					Title:   "Celebs",
					Content: "celebs",
				},
				models.Article{
					ID:      3,
					Title:   "Asmara",
					Content: "asmara",
				},
				models.Article{
					ID:      2,
					Title:   "Bola",
					Content: "bola",
				},
			},
		},
		testCase{
			Title:  "fetch-with-cursor",
			Num:    3,
			Cursor: "3",
			ExpectedResult: []models.Article{
				models.Article{
					ID:      2,
					Title:   "Bola",
					Content: "bola",
				},
				models.Article{
					ID:      1,
					Title:   "Tekno",
					Content: "tekno",
				},
			},
		},
	}

	for _, tc := range arrTestcase {
		m.T().Run(tc.Title, func(t *testing.T) {
			res, _, err := repo.Fetch(context.Background(), tc.Cursor, tc.Num)
			require.NoError(t, err)
			require.Equal(t, len(tc.ExpectedResult), len(res), tc.Title)
			for i, item := range res {
				assert.Equal(t, tc.ExpectedResult[i].ID, item.ID)
				assert.Equal(t, tc.ExpectedResult[i].Title, item.Title)
				assert.Equal(t, tc.ExpectedResult[i].Content, item.Content)
			}
		})
	}
}

func (m *mysqlArticleSuiteTest) TestGetByID() {
	// Prepare
	mockArticle := getMockArrArticle()[0]
	seedArticleData(m.T(), m.DBConn)
	repo := repoHandler.NewArticleRepository(m.DBConn)

	// Test the function
	res, err := repo.GetByID(context.Background(), mockArticle.ID)

	// Evaluate the results
	require.NoError(m.T(), err)
	assert.Equal(m.T(), mockArticle.ID, res.ID)
	assert.Equal(m.T(), mockArticle.Title, res.Title)
	assert.Equal(m.T(), mockArticle.Content, res.Content)
}

func (m *mysqlArticleSuiteTest) TestUpdate() {
	// Prepare
	mockArticle := getMockArrArticle()[0]
	seedArticleData(m.T(), m.DBConn)
	repo := repoHandler.NewArticleRepository(m.DBConn)
	mockArticle.UpdatedAt = time.Now()
	mockArticle.Title = "Teknologi" // previously only "tekno"

	// Test the function
	err := repo.Update(context.Background(), &mockArticle)

	// Evaluate the results
	require.NoError(m.T(), err)
	res := getArticleByID(m.T(), m.DBConn, mockArticle.ID)
	assert.NotNil(m.T(), res)
	assert.Equal(m.T(), mockArticle.ID, res.ID)
	assert.Equal(m.T(), mockArticle.Title, res.Title)
	assert.Equal(m.T(), mockArticle.Content, res.Content)
}

func (m *mysqlArticleSuiteTest) TestDelete() {
	// Prepare
	mockArticle := getMockArrArticle()[0]
	seedArticleData(m.T(), m.DBConn)
	repo := repoHandler.NewArticleRepository(m.DBConn)

	// Test the function
	err := repo.Delete(context.Background(), mockArticle.ID)

	// Evaluate the results
	require.NoError(m.T(), err)
	res := getArticleByID(m.T(), m.DBConn, mockArticle.ID)
	assert.Empty(m.T(), res.ID)    // because already deleted the article should be empty
	assert.Empty(m.T(), res.Title) // because already deleted the article should be empty
}
