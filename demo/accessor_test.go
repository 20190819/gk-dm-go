package demo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Age int
}

type useCase struct {
	name string

	// 输入
	entity interface{}
	field  string

	// 期望输出
	wantVal any
	wantErr error
}

func TestReflectAccessorGetFieldValue(t *testing.T) {
	testCases := []useCase{
		{
			name:    "nil",
			field:   "Age",
			wantVal: nil,
			wantErr: ErrInvalidEntity,
		},
		{
			name:    "struct",
			entity:  &User{Age: 18},
			field:   "Age",
			wantVal: 18,
			wantErr: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			accessor, err := NewReflectAccessor(testCase.entity)
			if err != nil {
				assert.Equal(t, testCase.wantErr, err)
				return
			}
			value, err := accessor.FieldValue(testCase.field)
			if err != nil {
				assert.Equal(t, testCase.wantErr, err)
				return
			}
			assert.Equal(t, testCase.wantVal, value)
		})
	}
}

func TestReflectAccessorSetFieldValue(t *testing.T) {
	testCases := []struct {
		name    string
		entity  *User
		field   string
		newVal  int
		wantErr error
	}{
		{
			name:    "normal",
			entity:  &User{Age: 18},
			field:   "Age",
			newVal:  28,
			wantErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			accessor, err := NewReflectAccessor(testCase.entity)
			if err != nil {
				assert.Equal(t, testCase.wantErr, err)
				return
			}
			err = accessor.setFieldValue(testCase.field, testCase.newVal)
			assert.Equal(t, testCase.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, testCase.newVal, testCase.entity.Age)
			fmt.Println("新值：", testCase.entity.Age)
		})
	}
}
