package demo

import (
	"errors"
	"reflect"
)

var ErrInvalidEntity = errors.New("invalid entity")
var ErrIllegal = errors.New("非法字段")
var ErrForbiddenUpdate = errors.New("无法设置新值的字段")

type ReflectAccessor struct {
	val reflect.Value
	typ reflect.Type
}

func NewReflectAccessor(val any) (*ReflectAccessor, error) {
	typ := reflect.TypeOf(val)
	if typ == nil || (typ.Kind() != reflect.Pointer) || (typ.Elem().Kind() != reflect.Struct) {
		return nil, ErrInvalidEntity
	}
	return &ReflectAccessor{reflect.ValueOf(val).Elem(), typ.Elem()}, nil
}

func (r *ReflectAccessor) FieldValue(field string) (int, error) {
	if _, ok := r.typ.FieldByName(field); !ok {
		return 0, ErrIllegal
	}
	return r.val.FieldByName(field).Interface().(int), nil
}

func (r *ReflectAccessor) setFieldValue(field string, val any) error {
	if _, ok := r.typ.FieldByName(field); !ok {
		return ErrIllegal
	}
	fdVal := r.val.FieldByName(field)
	if !fdVal.CanSet() {
		return ErrForbiddenUpdate
	}
	fdVal.Set(reflect.ValueOf(val))
	return nil
}
