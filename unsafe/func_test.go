package unsafe

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type args struct {
	val any
}

type useCaseItem struct {
	name    string
	args    args
	want    map[string]*funcInfo
	wantErr error
}

type Order struct {
	buyer  int64
	seller int64
}

func (o Order) GetBuyer() int64 {
	return o.buyer
}

type OrderV1 struct {
	buyer  int64
	seller int64
}

func (o *OrderV1) GetBuyer() int64 {
	return o.buyer
}

func TestIterateFunc(t *testing.T) {
	testCases := []useCaseItem{
		{
			name:    "type-nil",
			wantErr: errors.New("输入 nil"),
		},
		{
			name:    "type-base",
			args:    args{123},
			wantErr: errors.New("不支持的类型"),
		},
		{
			name: "type-struct",
			args: args{
				val: Order{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*funcInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					in:     []reflect.Type{reflect.TypeOf(Order{})},
					out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "type-struct-but-input-pointer",
			args: args{
				val: &Order{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*funcInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					in:     []reflect.Type{reflect.TypeOf(&Order{})},
					out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "type-ptr",
			args: args{
				val: &OrderV1{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*funcInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					in:     []reflect.Type{reflect.TypeOf(&OrderV1{})},
					out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "type-ptr-but-input-struct",
			args: args{
				val: OrderV1{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*funcInfo{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := iterateFunc(testCase.args.val)
			if err != nil {
				assert.Equal(t, testCase.wantErr, err)
				return
			}

			assert.Equal(t, testCase.want, got)
		})
	}
}
