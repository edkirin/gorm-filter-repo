package repository

import (
	"github.com/edkirin/gormfilterrepo/smartfilter"

	"gorm.io/gorm/schema"
)

type GetOptions struct {
	Only       *[]string
	RaiseError *bool
}

type GetMethod[T schema.Tabler] struct {
	repo *RepoBase[T]
}

func (m *GetMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m GetMethod[T]) Get(filter interface{}, options *GetOptions) (*T, error) {
	var (
		model T
	)

	query, err := smartfilter.ToQuery(model, filter, m.repo.dbConn)
	if err != nil {
		return nil, err
	}

	if options != nil {
		query = applyOptionOnly(query, options.Only)
	}

	result := query.First(&model)
	if result.Error == nil {
		return &model, nil
	}

	if options != nil && options.RaiseError != nil && *options.RaiseError {
		return nil, result.Error
	}
	return nil, nil
}
