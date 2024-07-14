package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMethod(t *testing.T) {
	t.Run("With filter", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		filter := MyModelFilter{
			Id: &id,
		}

		sql := "DELETE FROM my_models WHERE my_models.id = $1"
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 111))
		mock.ExpectCommit()

		deleted, err := repo.Delete(filter)
		assert.Equal(t, int64(111), deleted)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("With multiple filters", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id1 := uuid.New()
		id2 := uuid.New()
		id3 := uuid.New()
		cnt := 3456
		filter := MyModelFilter{
			Ids:   &[]uuid.UUID{id1, id2, id3},
			CntGT: &cnt,
		}

		sql := "DELETE FROM my_models WHERE my_models.id IN ($1,$2,$3) AND my_models.cnt > $4"
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id1, id2, id3, cnt).
			WillReturnResult(sqlmock.NewResult(1, 123))
		mock.ExpectCommit()

		deleted, err := repo.Delete(filter)
		assert.Equal(t, int64(123), deleted)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
