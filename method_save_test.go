package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveMethod(t *testing.T) {
	t.Run("Save new model", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		model := MyModel{
			Value: "some value",
			Cnt:   123,
		}

		sql := "INSERT INTO my_models (id,value,cnt) VALUES ($1,$2,$3)"
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(model.Id, model.Value, model.Cnt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := repo.Save(&model)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Update existing model", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		id := uuid.New()
		model := MyModel{
			Id:    &id,
			Value: "some value",
			Cnt:   123,
		}

		sql := "UPDATE my_models SET value=$1,cnt=$2 WHERE id = $3"
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(model.Value, model.Cnt, model.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := repo.Save(&model)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
