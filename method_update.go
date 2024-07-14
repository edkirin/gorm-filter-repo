package repository

import (
	"github.com/edkirin/gormfilterrepo/smartfilter"

	"gorm.io/gorm/schema"
)

type UpdateMethod[T schema.Tabler] struct {
	repo *RepoBase[T]
}

func (m *UpdateMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m UpdateMethod[T]) Update(filter interface{}, values map[string]any) (int64, error) {
	var (
		model T
	)

	query, err := smartfilter.ToQuery(model, filter, m.repo.dbConn)
	if err != nil {
		return 0, err
	}
	result := query.Model(&model).Updates(values)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
