package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type FieldAccessor interface {
	GetFieldVal(field string) (int, error)
	SetFieldVal(field string, val int) error
}

type FieldMeta struct {
	typ    reflect.Type
	offset uintptr
}

type UfAccessor struct {
	fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

var ErrInvalid = errors.New("invalid entity")

func NewUfAccessor(entity interface{}) (*UfAccessor, error) {
	if entity == nil {
		return nil, ErrInvalid
	}

	typ := reflect.TypeOf(entity)

	// 不是指针
	if typ.Kind() != reflect.Pointer {
		return nil, ErrInvalid
	}
	// 是指针，但不是结构体指针
	if typ.Elem().Kind() != reflect.Struct {
		return nil, ErrInvalid
	}

	typElem := typ.Elem()
	typElemNum := typElem.NumField()
	resFields := make(map[string]FieldMeta, typElemNum)
	for i := 0; i < typElemNum; i++ {
		fd := typElem.Field(i)
		resFields[fd.Name] = FieldMeta{
			typ:    fd.Type,
			offset: fd.Offset,
		}
	}

	return &UfAccessor{
		fields:     resFields,
		entityAddr: reflect.ValueOf(entity).UnsafePointer(),
	}, nil
}

func (u *UfAccessor) GetFieldVal(field string) (int, error) {

	if !u.validateField(field) {
		return 0, u.errNotFound()
	}

	fMeta := u.fields[field]
	ufp := unsafe.Pointer(uintptr(u.entityAddr) + fMeta.offset)
	// *(*int) 什么意思？
	return *(*int)(ufp), nil
}

func (u *UfAccessor) GetFieldAny(field string) (any, error) {

	if !u.validateField(field) {
		return 0, u.errNotFound()
	}

	fMeta := u.fields[field]
	res := reflect.NewAt(fMeta.typ, unsafe.Pointer(uintptr(u.entityAddr)+fMeta.offset))
	return res.Elem().Interface(), nil
}

func (u *UfAccessor) SetFieldVal(field string, val int) error {

	if !u.validateField(field) {
		return u.errNotFound()
	}

	fMeta := u.fields[field]
	ufp := unsafe.Pointer(uintptr(u.entityAddr) + fMeta.offset)
	*(*int)(ufp) = val
	return nil
}

func (u *UfAccessor) SetFieldValAny(field string, val int) error {

	if !u.validateField(field) {
		return u.errNotFound()
	}
	fMeta := u.fields[field]
	ufp := reflect.NewAt(fMeta.typ, unsafe.Pointer(uintptr(u.entityAddr)+fMeta.offset))
	if ufp.CanSet() {
		ufp.Set(reflect.ValueOf(val))
	}
	return nil
}

func (u *UfAccessor) validateField(field string) bool {
	_, ok := u.fields[field]
	return ok
}

func (u *UfAccessor) errNotFound() error {
	return errors.New("field not found")
}
