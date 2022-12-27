package unsafe

import (
	"errors"
	"fmt"
	"reflect"
)

type funcInfo struct {
	Name   string
	in     []reflect.Type
	out    []reflect.Type
	Result []any
}

func iterateFunc(val any) (map[string]*funcInfo, error) {

	if val == nil {
		return nil, errors.New("输入 nil")
	}

	typ := reflect.TypeOf(val)
	if !validateType(typ) {
		return nil, errors.New("不支持的类型")
	}

	numMethod := typ.NumMethod()
	res := make(map[string]*funcInfo, numMethod)

	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		mt := method.Type

		numIn := mt.NumIn()
		in := make([]reflect.Type, 0, numIn)
		for i := 0; i < numIn; i++ {
			in = append(in, mt.In(i))
		}

		numOut := mt.NumOut()
		out := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			out = append(out, mt.Out(j))
		}

		callRes := method.Func.Call([]reflect.Value{reflect.ValueOf(val)})
		retValues := make([]any, 0, len(callRes))
		for _, ct := range callRes {
			retValues = append(retValues, ct.Interface())
		}

		fmt.Println("name", method.Name)

		res[method.Name] = &funcInfo{
			Name:   method.Name,
			in:     in,
			out:    out,
			Result: retValues,
		}
	}

	return res, nil
}

func validateType(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct || (typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct)
}
