package repository

import (
	"github.com/edkirin/gormfilterrepo/smartfilter"

	"gorm.io/gorm/schema"
)

type DeleteMethod[T schema.Tabler] struct {
	repo *RepoBase[T]
}

func (m *DeleteMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m DeleteMethod[T]) Delete(filter interface{}) (int64, error) {
	var (
		model T
	)

	query, err := smartfilter.ToQuery(model, filter, m.repo.dbConn)
	if err != nil {
		return 0, err
	}
	result := query.Delete(&model)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
