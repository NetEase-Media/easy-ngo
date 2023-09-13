package xreflect

import (
	"reflect"

	"github.com/pkg/errors"
)

// In test if key in map or value in slice/array
func In(value interface{}, container interface{}) bool {
	containerValue := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if containerValue.Index(i).Interface() == value {
				return true
			}
		}
	case reflect.Map:
		if containerValue.MapIndex(reflect.ValueOf(value)).IsValid() {
			return true
		}
	default:
		return false
	}
	return false
}

// Override try override the fields of left object with right object's
func Override(left interface{}, right interface{}) error {
	if reflect.ValueOf(left).IsNil() || reflect.ValueOf(right).IsNil() {
		return errors.New("nil arg")
	}

	if reflect.ValueOf(left).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(right).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(left).Kind() != reflect.ValueOf(right).Kind() {
		return errors.New("must be a pointer and have same type")
	}

	oldVal := reflect.ValueOf(left).Elem()
	newVal := reflect.ValueOf(right).Elem()

	if oldVal.Type() != newVal.Type() {
		return errors.New("must have same type")
	}

	for i := 0; i < oldVal.NumField(); i++ {
		newField := newVal.Field(i)
		oldField := oldVal.Field(i)
		if oldField.Kind() == reflect.Struct && newField.Kind() == reflect.Struct {
			if err := Override(oldField.Addr().Interface(), newField.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if !newField.IsZero() {
			oldField.Set(newField)
		}
	}
	return nil
}
