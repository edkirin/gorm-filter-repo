package smartfilter

import (
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		WithoutQuotingCheck: true,
		Conn:                db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

type MyModel struct {
	Id    int
	Value string
}

func (m MyModel) TableName() string {
	return "my_models"
}

type HandleOperatorTestCase struct {
	name        string
	filterField FilterField
	expected    string
}

var (
	boolTrue    bool    = true
	boolFalse   bool    = false
	int64Value  int64   = -123456
	uint64Value uint64  = 123456
	floatValue  float64 = -123456.789
	strValue    string  = "Some Value"

	boolValues   = []bool{true, false}
	int64Values  = []int64{-123456, 1, 123456}
	uint64Values = []uint64{123456, 1234567, 1234568}
	floatValues  = []float64{-123456.789, -1, 123456.789}
	strValues    = []string{"First Value", "Second Value", "Third Value"}
)

func TestHandleOperatorEQ(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorEQ

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorEQ bool true",
			filterField: FilterField{
				Name:      "my_field",
				boolValue: &boolTrue,
				valueKind: reflect.Bool,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field = true ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorEQ bool false",
			filterField: FilterField{
				Name:      "my_field",
				boolValue: &boolFalse,
				valueKind: reflect.Bool,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field = false ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorEQ int64",
			filterField: FilterField{
				Name:      "my_field",
				intValue:  &int64Value,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field = -123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorEQ uint64",
			filterField: FilterField{
				Name:      "my_field",
				uintValue: &uint64Value,
				valueKind: reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field = 123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorEQ float",
			filterField: FilterField{
				Name:       "my_field",
				floatValue: &floatValue,
				valueKind:  reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field = -123456.789 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorEQ string",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field = 'Some Value' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorNE(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorNE

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorNE bool true",
			filterField: FilterField{
				Name:      "my_field",
				boolValue: &boolTrue,
				valueKind: reflect.Bool,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field != true ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNE bool false",
			filterField: FilterField{
				Name:      "my_field",
				boolValue: &boolFalse,
				valueKind: reflect.Bool,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field != false ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNE int64",
			filterField: FilterField{
				Name:      "my_field",
				intValue:  &int64Value,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field != -123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNE uint64",
			filterField: FilterField{
				Name:      "my_field",
				uintValue: &uint64Value,
				valueKind: reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field != 123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNE float",
			filterField: FilterField{
				Name:       "my_field",
				floatValue: &floatValue,
				valueKind:  reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field != -123456.789 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNE string",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field != 'Some Value' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorLIKE(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorLIKE

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorLIKE",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field LIKE '%Some Value%' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorILIKE(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorILIKE

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorILIKE",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field ILIKE '%Some Value%' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorGT(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorGT

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorGT int64",
			filterField: FilterField{
				Name:      "my_field",
				intValue:  &int64Value,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field > -123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorGT uint64",
			filterField: FilterField{
				Name:      "my_field",
				uintValue: &uint64Value,
				valueKind: reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field > 123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorGT float",
			filterField: FilterField{
				Name:       "my_field",
				floatValue: &floatValue,
				valueKind:  reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field > -123456.789 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorGT string",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field > 'Some Value' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorGE(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorGE

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorGE int64",
			filterField: FilterField{
				Name:      "my_field",
				intValue:  &int64Value,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field >= -123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorGE uint64",
			filterField: FilterField{
				Name:      "my_field",
				uintValue: &uint64Value,
				valueKind: reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field >= 123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorGE float",
			filterField: FilterField{
				Name:       "my_field",
				floatValue: &floatValue,
				valueKind:  reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field >= -123456.789 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorGE string",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field >= 'Some Value' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorLT(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorLT

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorLT int64",
			filterField: FilterField{
				Name:      "my_field",
				intValue:  &int64Value,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field < -123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorLT uint64",
			filterField: FilterField{
				Name:      "my_field",
				uintValue: &uint64Value,
				valueKind: reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field < 123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorLT float",
			filterField: FilterField{
				Name:       "my_field",
				floatValue: &floatValue,
				valueKind:  reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field < -123456.789 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorLT string",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field < 'Some Value' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorLE(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorLE

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorLE int64",
			filterField: FilterField{
				Name:      "my_field",
				intValue:  &int64Value,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field <= -123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorLE uint64",
			filterField: FilterField{
				Name:      "my_field",
				uintValue: &uint64Value,
				valueKind: reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field <= 123456 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorLE float",
			filterField: FilterField{
				Name:       "my_field",
				floatValue: &floatValue,
				valueKind:  reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field <= -123456.789 ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorLE string",
			filterField: FilterField{
				Name:      "my_field",
				strValue:  &strValue,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field <= 'Some Value' ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorIN(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorIN

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorIN bool",
			filterField: FilterField{
				Name:       "my_field",
				boolValues: &boolValues,
				valueKind:  reflect.Bool,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field IN (true,false) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorIN int64",
			filterField: FilterField{
				Name:      "my_field",
				intValues: &int64Values,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field IN (-123456,1,123456) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorIN uint64",
			filterField: FilterField{
				Name:       "my_field",
				uintValues: &uint64Values,
				valueKind:  reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field IN (123456,1234567,1234568) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorIN float",
			filterField: FilterField{
				Name:        "my_field",
				floatValues: &floatValues,
				valueKind:   reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field IN (-123456.789,-1,123456.789) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorIN string",
			filterField: FilterField{
				Name:      "my_field",
				strValues: &strValues,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field IN ('First Value','Second Value','Third Value') ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}

func TestHandleOperatorNOT_IN(t *testing.T) {
	db, _ := NewMockDB()
	testFunc := handleOperatorNOT_IN

	testCases := []HandleOperatorTestCase{
		{
			name: "handleOperatorNOT_IN bool",
			filterField: FilterField{
				Name:       "my_field",
				boolValues: &boolValues,
				valueKind:  reflect.Bool,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field NOT IN (true,false) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNOT_IN int64",
			filterField: FilterField{
				Name:      "my_field",
				intValues: &int64Values,
				valueKind: reflect.Int64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field NOT IN (-123456,1,123456) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNOT_IN uint64",
			filterField: FilterField{
				Name:       "my_field",
				uintValues: &uint64Values,
				valueKind:  reflect.Uint64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field NOT IN (123456,1234567,1234568) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNOT_IN float",
			filterField: FilterField{
				Name:        "my_field",
				floatValues: &floatValues,
				valueKind:   reflect.Float64,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field NOT IN (-123456.789,-1,123456.789) ORDER BY my_models.id LIMIT 1",
		},
		{
			name: "handleOperatorNOT_IN string",
			filterField: FilterField{
				Name:      "my_field",
				strValues: &strValues,
				valueKind: reflect.String,
			},
			expected: "SELECT * FROM my_models WHERE my_table.my_field NOT IN ('First Value','Second Value','Third Value') ORDER BY my_models.id LIMIT 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				query := tx.Model(&MyModel{})
				query = testFunc(query, "my_table", &testCase.filterField)
				return query.First(&MyModel{})
			})
			assert.Equal(t, testCase.expected, sql)
		})
	}
}
