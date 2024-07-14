package smartfilter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetFilterFields(t *testing.T) {
	type TestFilter struct {
		Alive         *bool      `filterfield:"alive,EQ"`
		Id            *int64     `filterfield:"id,EQ"`
		Ids           *[]uint    `filterfield:"id,IN"`
		IdsNot        *[]uint    `filterfield:"id,NOT_IN"`
		FirstName     *string    `filterfield:"first_name,EQ"`
		NotFirstName  *string    `filterfield:"first_name,NE"`
		FirstNameLike *string    `filterfield:"first_name,LIKE"`
		CreatedAt_GE  *time.Time `filterfield:"created_at,GE"`
		CreatedAt_GT  *time.Time `filterfield:"created_at,GT"`
		CreatedAt_LE  *time.Time `filterfield:"created_at,LE"`
		CreatedAt_LT  *time.Time `filterfield:"created_at,LT"`
	}

	t.Run("Parse filter fields", func(t *testing.T) {
		utc, _ := time.LoadLocation("UTC")

		var (
			alive         bool      = true
			id            int64     = 123456
			ids           []uint    = []uint{111, 222, 333, 444, 555}
			idsNot        []uint    = []uint{666, 777, 888, 999}
			firstName     string    = "Mirko"
			notFirstName  string    = "Pero"
			firstNameLike string    = "irko"
			createdTime   time.Time = time.Date(2024, 5, 26, 16, 8, 0, 0, utc)
		)
		filter := TestFilter{
			Alive:         &alive,
			Id:            &id,
			Ids:           &ids,
			IdsNot:        &idsNot,
			FirstName:     &firstName,
			NotFirstName:  &notFirstName,
			FirstNameLike: &firstNameLike,
			CreatedAt_GE:  &createdTime,
			CreatedAt_GT:  &createdTime,
			CreatedAt_LE:  &createdTime,
			CreatedAt_LT:  &createdTime,
		}
		result := getFilterFields(filter)

		fmt.Printf("%+v\n", result[0].value)
		fmt.Printf("%+v\n", &alive)

		assert.Equal(t, "Alive", result[0].name)
		assert.Equal(t, alive, result[0].value.Elem().Bool())
		assert.Equal(t, "alive,EQ", result[0].tagValue)

		assert.Equal(t, "Id", result[1].name)
		assert.Equal(t, id, result[1].value.Elem().Int())
		assert.Equal(t, "id,EQ", result[1].tagValue)

		assert.Equal(t, "Ids", result[2].name)
		assert.Equal(t, ids, result[2].value.Elem().Interface())
		assert.Equal(t, "id,IN", result[2].tagValue)

		assert.Equal(t, "IdsNot", result[3].name)
		assert.Equal(t, idsNot, result[3].value.Elem().Interface())
		assert.Equal(t, "id,NOT_IN", result[3].tagValue)

		assert.Equal(t, "FirstName", result[4].name)
		assert.Equal(t, firstName, result[4].value.Elem().String())
		assert.Equal(t, "first_name,EQ", result[4].tagValue)

		assert.Equal(t, "NotFirstName", result[5].name)
		assert.Equal(t, notFirstName, result[5].value.Elem().String())
		assert.Equal(t, "first_name,NE", result[5].tagValue)

		assert.Equal(t, "FirstNameLike", result[6].name)
		assert.Equal(t, firstNameLike, result[6].value.Elem().String())
		assert.Equal(t, "first_name,LIKE", result[6].tagValue)

		assert.Equal(t, "CreatedAt_GE", result[7].name)
		assert.Equal(t, createdTime, result[7].value.Elem().Interface())
		assert.Equal(t, "created_at,GE", result[7].tagValue)

		assert.Equal(t, "CreatedAt_GT", result[8].name)
		assert.Equal(t, createdTime, result[8].value.Elem().Interface())
		assert.Equal(t, "created_at,GT", result[8].tagValue)

		assert.Equal(t, "CreatedAt_LE", result[9].name)
		assert.Equal(t, createdTime, result[9].value.Elem().Interface())
		assert.Equal(t, "created_at,LE", result[9].tagValue)

		assert.Equal(t, "CreatedAt_LT", result[10].name)
		assert.Equal(t, createdTime, result[10].value.Elem().Interface())
		assert.Equal(t, "created_at,LT", result[10].tagValue)
	})

	t.Run("Skip nil fields", func(t *testing.T) {
		type TestFilter struct {
			Alive     *bool   `filterfield:"alive;EQ"`
			Id        *int64  `filterfield:"id;EQ"`
			Ids       *[]uint `filterfield:"id;IN"`
			IdsNot    *[]uint `filterfield:"id;NOT_IN"`
			FirstName *string `filterfield:"first_name;EQ"`
		}
		filter := TestFilter{}
		result := getFilterFields(filter)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Skip fields without filterfield tag", func(t *testing.T) {
		var (
			alive bool  = true
			id    int64 = 123456
		)
		type TestFilter struct {
			Alive *bool
			Id    *int64 `funnytag:"created_at;LT"`
		}
		filter := TestFilter{
			Alive: &alive,
			Id:    &id,
		}
		result := getFilterFields(filter)
		assert.Equal(t, 0, len(result))
	})
}

type TagParseTestCase struct {
	name     string
	tagValue string
	expected FilterField
}

func TestFilterField(t *testing.T) {
	testCases := []TagParseTestCase{
		{
			name:     "Parse without spaces",
			tagValue: "field=field_1;operator=EQ",
			expected: FilterField{
				Name:     "field_1",
				Operator: OperatorEQ,
			},
		},
		{
			name:     "Parse spaces between pairs",
			tagValue: "   field=field_2 ;   operator=LT   ",
			expected: FilterField{
				Name:     "field_2",
				Operator: OperatorLT,
			},
		},
		{
			name:     "Parse spaces between around keys and values",
			tagValue: "operator   = LIKE ; field = field_3",
			expected: FilterField{
				Name:     "field_3",
				Operator: OperatorLIKE,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			filterField, err := newFilterField(testCase.tagValue)
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected.Name, filterField.Name)
			assert.Equal(t, testCase.expected.Operator, filterField.Operator)
		})
	}

	t.Run("Fail on invalid tag value", func(t *testing.T) {
		filterField, err := newFilterField("field=field_1=fail; operator=EQ")
		assert.Nil(t, filterField)
		assert.EqualError(t, err, "invalid tag value: field=field_1=fail")
	})

	t.Run("Fail on invalid operator", func(t *testing.T) {
		filterField, err := newFilterField("field=field_1; operator=FAIL")
		assert.Nil(t, filterField)
		assert.EqualError(t, err, "unknown operator: FAIL")
	})

	t.Run("Fail on invalid value key", func(t *testing.T) {
		filterField, err := newFilterField("failkey=field_1; operator=FAIL")
		assert.Nil(t, filterField)
		assert.EqualError(t, err, "invalid value key: failkey")
	})

	t.Run("Fail on missing field name", func(t *testing.T) {
		filterField, err := newFilterField("operator=EQ")
		assert.Nil(t, filterField)
		assert.EqualError(t, err, "missing field name in tag: operator=EQ")
	})

	t.Run("Fail on missing operator", func(t *testing.T) {
		filterField, err := newFilterField("field=field_1")
		assert.Nil(t, filterField)
		assert.EqualError(t, err, "missing operator in tag: field=field_1")
	})
}

type filterWithoutQueryApplier struct{}

type filterWithQueryApplier struct{}

func (f filterWithQueryApplier) ApplyQuery(query *gorm.DB) *gorm.DB {
	return query
}

func TestSmartfilterApplyQuery(t *testing.T) {

	t.Run("Get query applier interface - without interface", func(t *testing.T) {
		f := filterWithoutQueryApplier{}
		queryApplier := getQueryApplierInterface(f)
		assert.Nil(t, queryApplier)
	})

	t.Run("Get query applier interface - with interface", func(t *testing.T) {
		f := filterWithQueryApplier{}
		queryApplier := getQueryApplierInterface(f)
		assert.NotNil(t, queryApplier)
	})
}
