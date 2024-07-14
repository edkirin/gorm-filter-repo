package smartfilter

import (
	"reflect"

	"gorm.io/gorm"
)

func handleOperatorEQ(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Bool:
		return applyFilterEQ(query, tableName, filterField, *filterField.boolValue)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterEQ(query, tableName, filterField, *filterField.intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterEQ(query, tableName, filterField, *filterField.uintValue)
	case reflect.Float32, reflect.Float64:
		return applyFilterEQ(query, tableName, filterField, *filterField.floatValue)
	case reflect.String:
		return applyFilterEQ(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorNE(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Bool:
		return applyFilterNE(query, tableName, filterField, *filterField.boolValue)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterNE(query, tableName, filterField, *filterField.intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterNE(query, tableName, filterField, *filterField.uintValue)
	case reflect.Float32, reflect.Float64:
		return applyFilterNE(query, tableName, filterField, *filterField.floatValue)
	case reflect.String:
		return applyFilterNE(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorLIKE(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.String:
		return applyFilterLIKE(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorILIKE(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.String:
		return applyFilterILIKE(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorGT(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterGT(query, tableName, filterField, *filterField.intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterGT(query, tableName, filterField, *filterField.uintValue)
	case reflect.Float32, reflect.Float64:
		return applyFilterGT(query, tableName, filterField, *filterField.floatValue)
	case reflect.String:
		return applyFilterGT(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorGE(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterGE(query, tableName, filterField, *filterField.intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterGE(query, tableName, filterField, *filterField.uintValue)
	case reflect.Float32, reflect.Float64:
		return applyFilterGE(query, tableName, filterField, *filterField.floatValue)
	case reflect.String:
		return applyFilterGE(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorLT(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterLT(query, tableName, filterField, *filterField.intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterLT(query, tableName, filterField, *filterField.uintValue)
	case reflect.Float32, reflect.Float64:
		return applyFilterLT(query, tableName, filterField, *filterField.floatValue)
	case reflect.String:
		return applyFilterLT(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorLE(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterLE(query, tableName, filterField, *filterField.intValue)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterLE(query, tableName, filterField, *filterField.uintValue)
	case reflect.Float32, reflect.Float64:
		return applyFilterLE(query, tableName, filterField, *filterField.floatValue)
	case reflect.String:
		return applyFilterLE(query, tableName, filterField, *filterField.strValue)
	}
	return nil
}

func handleOperatorIN(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Bool:
		return applyFilterIN(query, tableName, filterField, filterField.boolValues)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterIN(query, tableName, filterField, filterField.intValues)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterIN(query, tableName, filterField, filterField.uintValues)
	case reflect.Float32, reflect.Float64:
		return applyFilterIN(query, tableName, filterField, filterField.floatValues)
	case reflect.String:
		return applyFilterIN(query, tableName, filterField, filterField.strValues)
	}
	return nil
}

func handleOperatorNOT_IN(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB {
	switch filterField.valueKind {
	case reflect.Bool:
		return applyFilterNOT_IN(query, tableName, filterField, filterField.boolValues)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return applyFilterNOT_IN(query, tableName, filterField, filterField.intValues)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return applyFilterNOT_IN(query, tableName, filterField, filterField.uintValues)
	case reflect.Float32, reflect.Float64:
		return applyFilterNOT_IN(query, tableName, filterField, filterField.floatValues)
	case reflect.String:
		return applyFilterNOT_IN(query, tableName, filterField, filterField.strValues)
	}
	return nil
}
