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

func TestUfAccessor_GetFieldVal(t *testing.T) {
	tests := []testCaseItemUnSafe{
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

	for _, tc := range tests {
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
