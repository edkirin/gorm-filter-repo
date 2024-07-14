package repository

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		WithoutQuotingCheck: true,
		Conn:                sqldb,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return sqldb, gormdb, mock
}

type MyModel struct {
	Id    *uuid.UUID `gorm:"type(uuid);unique"`
	Value string
	Cnt   int
}

func (m MyModel) TableName() string {
	return "my_models"
}

type MyModelFilter struct {
	Id    *uuid.UUID   `filterfield:"field=id;operator=EQ"`
	Ids   *[]uuid.UUID `filterfield:"field=id;operator=IN"`
	Value *string      `filterfield:"field=value;operator=EQ"`
	CntGT *int         `filterfield:"field=cnt;operator=GT"`
}

func TestListMethod(t *testing.T) {
	t.Run("With ordering", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		filter := MyModelFilter{}
		options := ListOptions{
			Ordering: &[]Order{
				{
					Field:     "id",
					Direction: OrderASC,
				},
				{
					Field:     "cnt",
					Direction: OrderDESC,
				},
			},
		}

		sql := "SELECT * FROM my_models ORDER BY id,cnt DESC"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql)))

		_, err := repo.List(filter, &options)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("With limit", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		filter := MyModelFilter{}
		options := ListOptions{
			Pagination: &Pagination{
				Limit: 111,
			},
		}

		sql := "SELECT * FROM my_models LIMIT $1"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(options.Pagination.Limit)

		_, err := repo.List(filter, &options)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("With offset", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		filter := MyModelFilter{}
		options := ListOptions{
			Pagination: &Pagination{
				Offset: 222,
			},
		}

		sql := "SELECT * FROM my_models OFFSET $1"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(options.Pagination.Offset)

		_, err := repo.List(filter, &options)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("With limit and offset", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		filter := MyModelFilter{}
		options := ListOptions{
			Pagination: &Pagination{
				Limit:  111,
				Offset: 222,
			},
		}

		sql := "SELECT * FROM my_models LIMIT $1 OFFSET $2"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(options.Pagination.Limit, options.Pagination.Offset)

		_, err := repo.List(filter, &options)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Simple filter", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		filter := MyModelFilter{
			Id: &id,
		}

		sql := "SELECT * FROM my_models WHERE my_models.id = $1"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id)

		_, err := repo.List(filter, nil)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Multiple filter values", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		count := 456
		value := "some value"
		filter := MyModelFilter{
			Id:    &id,
			Value: &value,
			CntGT: &count,
		}

		sql := "SELECT * FROM my_models WHERE my_models.id = $1 AND my_models.value = $2 AND my_models.cnt > $3"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id, value, count)

		_, err := repo.List(filter, nil)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Multiple filter values and pagination", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		count := 456
		value := "some value"
		filter := MyModelFilter{
			Id:    &id,
			Value: &value,
			CntGT: &count,
		}
		options := ListOptions{
			Pagination: &Pagination{
				Offset: 111,
				Limit:  222,
			},
		}

		sql := "SELECT * FROM my_models WHERE my_models.id = $1 AND my_models.value = $2 AND my_models.cnt > $3 LIMIT $4 OFFSET $5"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id, value, count, options.Pagination.Limit, options.Pagination.Offset)

		_, err := repo.List(filter, &options)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Only id and count", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		filter := MyModelFilter{}
		options := ListOptions{
			Only: &[]string{
				"id",
				"cnt",
			},
		}

		sql := "SELECT id,cnt FROM my_models"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql)))

		_, err := repo.List(filter, &options)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
