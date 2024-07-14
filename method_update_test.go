package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMethod(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		sqldb, db, mock := NewMockDB()
		defer sqldb.Close()

		repo := RepoBase[MyModel]{}
		repo.Init(db, nil)

		cnt := 10
		filter := MyModelFilter{
			CntGT: &cnt,
		}
		values := map[string]any{
			"cnt":   111,
			"value": 222,
		}

		sql := "UPDATE my_models SET cnt=$1,value=$2 WHERE my_models.cnt > $3"
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf("^%s$", regexp.QuoteMeta(sql))).
			WithArgs(values["cnt"], values["value"], cnt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := repo.Update(filter, values)
		assert.Nil(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
