package repository

import (
	"github.com/edkirin/gormfilterrepo/smartfilter"

	"gorm.io/gorm/schema"
)

type ListOptions struct {
	Only       []string
	Ordering   []Order
	Pagination *Pagination
	Joins      []string
}

type ListMethod[T schema.Tabler] struct {
	repo *RepoBase[T]
}

func (m *ListMethod[T]) Init(repo *RepoBase[T]) {
	m.repo = repo
}

func (m ListMethod[T]) List(filter interface{}, options *ListOptions) (*[]T, error) {
	var (
		model  T
		models []T
	)

	query, err := smartfilter.ToQuery(model, filter, m.repo.dbConn)
	if err != nil {
		return nil, err
	}

	if options != nil {
		query = applyJoins(query, options.Joins)
		query = applyOptionOnly(query, options.Only)
		query = applyOptionOrdering(query, options.Ordering)
		query = applyOptionPagination(query, options.Pagination)
	}

	query.Find(&models)
	return &models, nil
}
