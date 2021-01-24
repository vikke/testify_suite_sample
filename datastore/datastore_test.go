package datastore

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/khaiql/dbcleaner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

const dataSourceName = "./database.sqlite3"

type DatabaseTestSuite struct {
	suite.Suite

	db      *sqlx.DB
	cleaner dbcleaner.DbCleaner
}

/*
type PostRepositoryTestSuite struct {
	DatabaseTestSuite
}
*/

func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) setupDB() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", dataSourceName)
}

func (s *DatabaseTestSuite) setupDBCleaner() dbcleaner.DbCleaner {
	cleaner := dbcleaner.New()
	sqlite := engine.NewSqliteEngine(dataSourceName)
	cleaner.SetEngine(sqlite)

	return cleaner
}

// スイートのセットアップ時にDBとクリーナーをセットする
func (s *DatabaseTestSuite) SetupSuite() {
	s.T().Log("# SetupSuite")
	db, err := s.setupDB()
	s.Require().NoError(err)
	s.db = db
	s.cleaner = s.setupDBCleaner()
	db.Exec("drop database posts")
	db.Exec("create table posts (id integer, title text, body text)")
}

// テストのセットアップ時にDBに繋いでデータを空にしておく
func (s *DatabaseTestSuite) SetupTest() {
	s.T().Log("## SetupTest")
	s.cleanData()
	s.insertData()
}

// テスト実行後にデータを掃除
func (s *DatabaseTestSuite) TearDownTest() {
	s.T().Log("## TearDownTest")
	s.cleanData()
}

// クリーナーとDBを閉じる
func (s *DatabaseTestSuite) TearDownSuite() {
	s.T().Log("# TearDownSuite")
	s.cleaner.Close()
	s.db.Close()
}

func (s *DatabaseTestSuite) insertData() {
	posts := []Post{
		Post{ID: 1, Title: "title1", Body: "body1"},
		Post{ID: 2, Title: "title2", Body: "body2"},
		Post{ID: 3, Title: "title3", Body: "body3"},
	}
	for _, p := range posts {
		_, err := s.db.Exec("insert into posts (id, title, body) values (?, ?, ?)", p.ID, p.Title, p.Body)
		s.Require().NoError(err)
	}
}

func (s *DatabaseTestSuite) cleanData() {
	s.cleaner.Clean("posts")
}

func (s *DatabaseTestSuite) TestPostRepositoryGetByID() {
	s.T().Log("##### TestPostRepositoryGetByID")
	repo := NewPostRepository(s.db)
	got, err := repo.GetByID(2)

	s.Assert().NoError(err)
	s.Assert().Equal(2, got.ID)
	s.Assert().Equal("title2", got.Title)
	s.Assert().Equal("body2", got.Body)
}

func (s *DatabaseTestSuite) TestPostRepositoryList() {
	s.T().Log("##### TestPostRepositoryList")
	repo := NewPostRepository(s.db)
	got, err := repo.List()

	s.Assert().NoError(err)
	s.Assert().Equal(3, len(got))
	s.Assert().Equal(3, got[0].ID)
	s.Assert().Equal(2, got[1].ID)
}
