package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExistsMethod(t *testing.T) {
	t.Run("Without repo options", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		filter := MyModelFilter{
			Id: &id,
		}

		sql := "SELECT id FROM my_models WHERE my_models.id = $1 LIMIT $2"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).WithArgs(id, 1)

		result, err := repo.Exists(filter)
		assert.False(t, result)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("With id field set in repo options", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, &RepoOptions{IdField: "some_other_pk"})

		id := uuid.New()
		filter := MyModelFilter{
			Id: &id,
		}

		sql := "SELECT some_other_pk FROM my_models WHERE my_models.id = $1 LIMIT $2"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).WithArgs(id, 1)

		result, err := repo.Exists(filter)
		assert.False(t, result)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
