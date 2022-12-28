package unsafe

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCaseItemUnSafe struct {
	name    string
	entity  interface{}
	field   string
	wantVal int
	wantErr error
}

type User struct {
	Age int
}

func TestUfAccessorGetFieldVal(t *testing.T) {

	testCases := []testCaseItemUnSafe{
		{
			name:    "invalid-field",
			entity:  &User{Age: 18},
			field:   "Age1",
			wantErr: errors.New("field not found"),
		},
		{
			name:    "normal-field",
			entity:  &User{Age: 18},
			field:   "Age",
			wantVal: 18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor, err := NewUfAccessor(tc.entity)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			val, err := accessor.GetFieldVal(tc.field)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

type ufAccessorTestCaseItem struct {
	name    string
	entity  *User
	field   string
	newVal  int
	wantErr error
}

func TestUnsafeAccessorSetFieldVal(t *testing.T) {

	testCases := []ufAccessorTestCaseItem{
		{
			name:    "not found case",
			entity:  &User{},
			field:   "Age1",
			newVal:  20,
			wantErr: errors.New("field not found"),
		},
		{
			name:   "normal case",
			entity: &User{},
			field:  "Age",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor, err := NewUfAccessor(tc.entity)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			err = accessor.SetFieldVal(tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.newVal, tc.entity.Age)
		})
	}
}

func TestUnsafeAccessorSetFieldValAny(t *testing.T) {
	testCases := []ufAccessorTestCaseItem{
		{
			name:   "normal case",
			entity: &User{},
			field:  "Age",
			newVal: 18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor, err := NewUfAccessor(tc.entity)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			err = accessor.SetFieldValAny(tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.newVal, tc.entity.Age)
		})
	}
}
