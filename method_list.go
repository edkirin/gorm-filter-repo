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
		query = ApplyJoins(query, options.Joins)
		query = ApplyOptionOnly(query, options.Only)
		query = ApplyOptionOrdering(query, options.Ordering)
		query = ApplyOptionPagination(query, options.Pagination)
	}

	query.Find(&models)
	return &models, nil
}
