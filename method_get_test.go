package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMethod(t *testing.T) {
	t.Run("With id filter, suppress errors", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		filter := MyModelFilter{
			Id: &id,
		}
		options := GetOptions{}

		sql := "SELECT * FROM my_models WHERE my_models.id = $1 ORDER BY my_models.id LIMIT $2"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).WithArgs(id, 1)

		result, err := repo.Get(filter, &options)
		assert.Nil(t, result)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("With id filter, raise error", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		filter := MyModelFilter{
			Id: &id,
		}
		options := GetOptions{
			RaiseError: true,
		}

		sql := "SELECT * FROM my_models WHERE my_models.id = $1 ORDER BY my_models.id LIMIT $2"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id, 1).
			WillReturnError(nil)

		result, err := repo.Get(filter, &options)
		assert.Nil(t, result)
		assert.NotNil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
