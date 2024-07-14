package smartfilter

import (
	"fmt"

	"gorm.io/gorm"
)

func applyFilterEQ[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s = ?", tableName, filterField.Name), value)
}

func applyFilterNE[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s != ?", tableName, filterField.Name), value)
}

func applyFilterLIKE(query *gorm.DB, tableName string, filterField *FilterField, value string) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s LIKE ?", tableName, filterField.Name), fmt.Sprintf("%%%s%%", value))
}

func applyFilterILIKE(query *gorm.DB, tableName string, filterField *FilterField, value string) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s ILIKE ?", tableName, filterField.Name), fmt.Sprintf("%%%s%%", value))
}

func applyFilterGT[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s > ?", tableName, filterField.Name), value)
}

func applyFilterGE[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s >= ?", tableName, filterField.Name), value)
}

func applyFilterLT[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s < ?", tableName, filterField.Name), value)
}

func applyFilterLE[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s <= ?", tableName, filterField.Name), value)
}

func applyFilterIN[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value *[]T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s IN (?)", tableName, filterField.Name), *value)
}

func applyFilterNOT_IN[T bool | int64 | uint64 | float64 | string](
	query *gorm.DB, tableName string, filterField *FilterField, value *[]T,
) *gorm.DB {
	return query.Where(fmt.Sprintf("%s.%s NOT IN (?)", tableName, filterField.Name), *value)
}
