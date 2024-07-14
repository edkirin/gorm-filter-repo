package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCountMethod(t *testing.T) {
	t.Run("Count without filter", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		filter := MyModelFilter{}

		sql := "SELECT count(*) FROM my_models"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql)))

		result, err := repo.Count(filter)
		assert.Equal(t, result, int64(0))
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Count with filter", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id1 := uuid.New()
		id2 := uuid.New()
		id3 := uuid.New()
		filter := MyModelFilter{
			Ids: &[]uuid.UUID{id1, id2, id3},
		}

		sql := "SELECT count(*) FROM my_models WHERE my_models.id IN ($1,$2,$3)"
		mock.ExpectQuery(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(id1, id2, id3)

		result, err := repo.Count(filter)
		assert.Equal(t, result, int64(0))
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
