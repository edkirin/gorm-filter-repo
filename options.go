package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type Pagination struct {
	Offset int
	Limit  int
}

type OrderDirection string

const (
	OrderASC  OrderDirection = "ASC"
	OrderDESC OrderDirection = "DESC"
)

type Order struct {
	Field     string
	Direction OrderDirection
}

func applyOptionOnly(query *gorm.DB, only *[]string) *gorm.DB {
	if only == nil || len(*only) == 0 {
		return query
	}
	query = query.Select(*only)
	return query
}

func applyOptionOrdering(query *gorm.DB, ordering *[]Order) *gorm.DB {
	if ordering == nil || len(*ordering) == 0 {
		return query
	}

	for _, order := range *ordering {
		if len(order.Direction) == 0 || order.Direction == OrderASC {
			query = query.Order(order.Field)
		} else {
			query = query.Order(fmt.Sprintf("%s %s", order.Field, order.Direction))
		}
	}
	return query
}

func applyOptionPagination(query *gorm.DB, pagination *Pagination) *gorm.DB {
	if pagination == nil {
		return query
	}

	if pagination.Limit != 0 {
		query = query.Limit(pagination.Limit)
	}
	if pagination.Offset != 0 {
		query = query.Offset(pagination.Offset)
	}
	return query
}
