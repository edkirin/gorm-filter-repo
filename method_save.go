package repository

import (
	"gorm.io/gorm/schema"
)

type SaveMethod[T schema.Tabler] struct {
	repo     *RepoBase[T]
	PreSave  func(model *T) error
	PostSave func(model *T) error
}

func (m *SaveMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m SaveMethod[T]) Save(model *T) (*T, error) {
	if m.PreSave != nil {
		err := m.PreSave(model)
		if err != nil {
			return nil, err
		}
	}

	result := m.repo.dbConn.Save(model)

	if result.Error != nil && m.PostSave != nil {
		err := m.PostSave(model)
		if err != nil {
			return nil, err
		}
	}

	return model, result.Error
}
