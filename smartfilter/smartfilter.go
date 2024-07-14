package smartfilter

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const TAG_NAME = "filterfield"
const TAG_PAIRS_SEPARATOR = ";"
const TAG_LIST_SEPARATOR = ","
const TAG_KEYVALUE_SEPARATOR = "="

type handlerFunc func(query *gorm.DB, tableName string, filterField *FilterField) *gorm.DB

type QueryApplier interface {
	ApplyQuery(query *gorm.DB) *gorm.DB
}

var operatorHandlers = map[Operator]handlerFunc{
	OperatorEQ:     handleOperatorEQ,
	OperatorNE:     handleOperatorNE,
	OperatorGT:     handleOperatorGT,
	OperatorGE:     handleOperatorGE,
	OperatorLT:     handleOperatorLT,
	OperatorLE:     handleOperatorLE,
	OperatorLIKE:   handleOperatorLIKE,
	OperatorILIKE:  handleOperatorILIKE,
	OperatorIN:     handleOperatorIN,
	OperatorNOT_IN: handleOperatorNOT_IN,
}

type ReflectedStructField struct {
	name     string
	value    reflect.Value
	tagValue string
}

func getFilterFields(filter interface{}) []ReflectedStructField {
	res := make([]ReflectedStructField, 0)

	st := reflect.TypeOf(filter)
	reflectValue := reflect.ValueOf(filter)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tagValue := field.Tag.Get(TAG_NAME)

		// skip field if filter tag is not present
		if len(tagValue) == 0 {
			continue
		}

		// get field value
		fieldValue := reflectValue.FieldByName(field.Name)

		// skip field if value is nil
		if fieldValue.IsNil() {
			continue
		}

		res = append(res, ReflectedStructField{
			name:     field.Name,
			tagValue: tagValue,
			value:    fieldValue,
		})
	}
	return res
}

func getQueryApplierInterface(filter interface{}) QueryApplier {
	queryApplier, ok := filter.(QueryApplier)
	if ok {
		return queryApplier
	}
	return nil
}

func ToQuery(model schema.Tabler, filter interface{}, query *gorm.DB) (*gorm.DB, error) {
	st := reflect.TypeOf(filter)

	tableName := model.TableName()
	modelName := st.Name()

	fields := getFilterFields(filter)
	for _, field := range fields {
		filterField, err := newFilterField(field.tagValue)
		if err != nil {
			return nil, fmt.Errorf("%s.%s: %s", modelName, field.name, err)
		}

		// must be called!
		filterField.setValueFromReflection(field.value)

		operatorHandler, ok := operatorHandlers[filterField.Operator]
		if !ok {
			return nil, fmt.Errorf("no handler for operator %s", filterField.Operator)
		}

		query = operatorHandler(query, tableName, filterField)
		if query == nil {
			return nil, fmt.Errorf("invalid field type for operator %s", filterField.Operator)
		}
	}

	// apply custom filters, if interface exists
	queryApplier := getQueryApplierInterface(filter)
	if queryApplier != nil {
		query = queryApplier.ApplyQuery(query)
	}

	return query, nil
}

func splitTrim(value string, separator string) []string {
	var out []string = []string{}
	for _, s := range strings.Split(value, separator) {
		if len(s) == 0 {
			continue
		}
		out = append(out, strings.TrimSpace(s))
	}
	return out
}

func newFilterField(tagValue string) (*FilterField, error) {
	filterField := FilterField{}

	for _, pair := range splitTrim(tagValue, TAG_PAIRS_SEPARATOR) {
		kvs := splitTrim(pair, TAG_KEYVALUE_SEPARATOR)
		if len(kvs) != 2 {
			return nil, fmt.Errorf("invalid tag value: %s", strings.TrimSpace(pair))
		}
		key := kvs[0]
		value := kvs[1]

		switch key {
		case "field":
			filterField.Name = value
		case "operator":
			operator := Operator(value)
			if !slices.Contains(OPERATORS, operator) {
				return nil, fmt.Errorf("unknown operator: %s", operator)
			}
			filterField.Operator = operator
		default:
			return nil, fmt.Errorf("invalid value key: %s", key)
		}
	}

	if len(filterField.Name) == 0 {
		return nil, fmt.Errorf("missing field name in tag: %s", tagValue)
	}
	if len(filterField.Operator) == 0 {
		return nil, fmt.Errorf("missing operator in tag: %s", tagValue)
	}

	return &filterField, nil
}
