package helper_test

import (
	"reflect"
	"testing"

	"github.com/cpd007/myredis/helper"
)

func TestToStringArray(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		d       any
		want    []string
		wantErr bool
	}{
		{
			name: "converts strings",
			d:    []any{"one", "two", "three"},
			want: []string{"one", "two", "three"},
		},
		{
			name: "accepts empty array",
			d:    []any{},
			want: []string{},
		},
		{
			name:    "rejects non-array input",
			d:       "value",
			wantErr: true,
		},
		{
			name:    "rejects typed string array",
			d:       []string{"one", "two"},
			wantErr: true,
		},
		{
			name:    "rejects non-string element",
			d:       []any{"one", 2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := helper.ToStringArray(tt.d)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToStringArray() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ToStringArray() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
