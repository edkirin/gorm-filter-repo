package repository

import (
	"gorm.io/gorm/schema"
)

type SaveMethod[T schema.Tabler] struct {
	repo *RepoBase[T]
}

func (m *SaveMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m SaveMethod[T]) Save(model *T) (*T, error) {
	result := m.repo.dbConn.Save(model)
	return model, result.Error
}
