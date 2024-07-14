package repository

import (
	"github.com/edkirin/gormfilterrepo/smartfilter"

	"gorm.io/gorm/schema"
)

type CountMethod[T schema.Tabler] struct {
	repo *RepoBase[T]
}

func (m *CountMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m CountMethod[T]) Count(filter interface{}) (int64, error) {
	var (
		model T
		count int64
	)

	query := m.repo.dbConn.Model(model)

	query, err := smartfilter.ToQuery(model, filter, query)
	if err != nil {
		return 0, err
	}

	query.Count(&count)
	return count, nil
}
