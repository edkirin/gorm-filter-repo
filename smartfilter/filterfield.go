package smartfilter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type FilterField struct {
	Name     string
	Operator Operator

	valueKind   reflect.Kind
	boolValue   *bool
	intValue    *int64
	uintValue   *uint64
	floatValue  *float64
	strValue    *string
	boolValues  *[]bool
	intValues   *[]int64
	uintValues  *[]uint64
	floatValues *[]float64
	strValues   *[]string
}

func (ff *FilterField) setValueFromReflection(v reflect.Value) {
	fn := typeGetter(v.Type())
	fn(ff, v)
}

func (ff *FilterField) appendStr(value string) {
	var valueArray []string

	if ff.strValues == nil {
		valueArray = make([]string, 0)
	} else {
		valueArray = *ff.strValues
	}
	valueArray = append(valueArray, value)
	ff.strValues = &valueArray
	ff.valueKind = reflect.String
}

func (ff *FilterField) appendBool(value bool) {
	var valueArray []bool

	if ff.boolValues == nil {
		valueArray = make([]bool, 0)
	} else {
		valueArray = *ff.boolValues
	}
	valueArray = append(valueArray, value)
	ff.boolValues = &valueArray
	ff.valueKind = reflect.Bool
}

func (ff *FilterField) appendInt(value int64) {
	var valueArray []int64

	if ff.boolValues == nil {
		valueArray = make([]int64, 0)
	} else {
		valueArray = *ff.intValues
	}
	valueArray = append(valueArray, value)
	ff.intValues = &valueArray
	ff.valueKind = reflect.Int
}

func (ff *FilterField) appendUint(value uint64) {
	var valueArray []uint64

	if ff.boolValues == nil {
		valueArray = make([]uint64, 0)
	} else {
		valueArray = *ff.uintValues
	}
	valueArray = append(valueArray, value)
	ff.uintValues = &valueArray
	ff.valueKind = reflect.Int
}

func (ff *FilterField) appendFloat(value float64) {
	var valueArray []float64

	if ff.boolValues == nil {
		valueArray = make([]float64, 0)
	} else {
		valueArray = *ff.floatValues
	}
	valueArray = append(valueArray, value)
	ff.floatValues = &valueArray
	ff.valueKind = reflect.Int
}

type valueGetterFunc func(ff *FilterField, v reflect.Value) error

func boolValueGetter(ff *FilterField, v reflect.Value) error {
	value := v.Bool()
	ff.boolValue = &value
	ff.valueKind = reflect.Bool
	return nil
}

func intValueGetter(ff *FilterField, v reflect.Value) error {
	value := v.Int()
	ff.intValue = &value
	ff.valueKind = reflect.Int
	return nil
}

func uintValueGetter(ff *FilterField, v reflect.Value) error {
	value := v.Uint()
	ff.uintValue = &value
	ff.valueKind = reflect.Uint
	return nil
}

func floatValueGetter(ff *FilterField, v reflect.Value) error {
	value := v.Float()
	ff.floatValue = &value
	ff.valueKind = reflect.Float64
	return nil
}

func strValueGetter(ff *FilterField, v reflect.Value) error {
	value := v.String()
	ff.strValue = &value
	ff.valueKind = reflect.String
	return nil
}

func timeValueGetter(ff *FilterField, v reflect.Value) error {
	value, err := timeValueToStr(v)
	if err != nil {
		return err
	}
	ff.strValue = &value
	ff.valueKind = reflect.String
	return nil
}

func uuidValueGetter(ff *FilterField, v reflect.Value) error {
	value, err := uuidValueToStr(v)
	if err != nil {
		return err
	}
	ff.strValue = &value
	ff.valueKind = reflect.String
	return nil
}

func timeValueToStr(v reflect.Value) (string, error) {
	value, ok := v.Interface().(time.Time)
	if !ok {
		return "", fmt.Errorf("error converting interface to time")
	}
	return value.Format(time.RFC3339), nil
}

func uuidValueToStr(v reflect.Value) (string, error) {
	value, ok := v.Interface().(uuid.UUID)
	if !ok {
		return "", fmt.Errorf("error converting interface to uuid")
	}
	return value.String(), nil
}

func unsupportedValueGetter(ff *FilterField, v reflect.Value) error {
	return fmt.Errorf("unsupported type: %v", v.Type())
}

func typeGetter(t reflect.Type) valueGetterFunc {
	return newTypeGetter(t, true)
}

func newTypeGetter(t reflect.Type, allowAddr bool) valueGetterFunc {
	// If we have a non-pointer value whose type implements
	// Marshaler with a value receiver, then we're better off taking
	// the address of the value - otherwise we end up with an
	// allocation as we cast the value to an interface.
	// if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(marshalerType) {
	// 	return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
	// }
	// if t.Implements(marshalerType) {
	// 	return marshalerEncoder
	// }
	// if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(textMarshalerType) {
	// 	return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
	// }
	// if t.Implements(textMarshalerType) {
	// 	return textMarshalerEncoder
	// }

	switch t.Kind() {
	case reflect.Bool:
		return boolValueGetter
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intValueGetter
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintValueGetter
	case reflect.Float32, reflect.Float64:
		return floatValueGetter
	case reflect.String:
		return strValueGetter
	// case reflect.Interface:
	// 	return interfaceEncoder
	case reflect.Struct:
		// check if value type is time
		if t == reflect.TypeOf(time.Time{}) {
			return timeValueGetter
		}
	// case reflect.Map:
	// 	return newMapEncoder(t)
	case reflect.Slice:
		return newSliceGetter(t)
	case reflect.Array:
		if t == reflect.TypeOf(uuid.UUID{}) {
			return uuidValueGetter
		}
		return newArrayGetter(t)
	case reflect.Pointer:
		return newPtrValueGetter(t)
	}
	return unsupportedValueGetter
}

type ptrValueGetter struct {
	elemGetter valueGetterFunc
}

func (pvg ptrValueGetter) getValue(ff *FilterField, v reflect.Value) error {
	pvg.elemGetter(ff, v.Elem())
	return nil
}

func newPtrValueGetter(t reflect.Type) valueGetterFunc {
	enc := ptrValueGetter{elemGetter: typeGetter(t.Elem())}
	return enc.getValue
}

type arrayGetter struct {
	elemGetter valueGetterFunc
}

func (ag arrayGetter) getValue(ff *FilterField, v reflect.Value) error {
	ag.elemGetter(ff, v.Elem())
	return nil
}

func newArrayGetter(t reflect.Type) valueGetterFunc {
	enc := arrayGetter{elemGetter: typeGetter(t.Elem())}
	return enc.getValue
}

type sliceGetter struct {
	elemGetter valueGetterFunc
}

func (sg sliceGetter) getValue(ff *FilterField, v reflect.Value) error {
	for n := range v.Len() {
		element := v.Index(n)

		switch element.Kind() {
		case reflect.Bool:
			ff.appendBool(element.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ff.appendInt(element.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			ff.appendUint(element.Uint())
		case reflect.Float32, reflect.Float64:
			ff.appendFloat(element.Float())
		case reflect.String:
			ff.appendStr(element.String())
		case reflect.Struct:
			if element.Type() == reflect.TypeOf(time.Time{}) {
				value, err := timeValueToStr(element)
				if err != nil {
					return err
				}
				ff.appendStr(value)
			}
		case reflect.Array:
			if element.Type() == reflect.TypeOf(uuid.UUID{}) {
				value, err := uuidValueToStr(element)
				if err != nil {
					return err
				}
				ff.appendStr(value)
			}
		}
	}
	return nil
}

func newSliceGetter(t reflect.Type) valueGetterFunc {
	enc := sliceGetter{elemGetter: newArrayGetter(t)}
	return enc.getValue
}
