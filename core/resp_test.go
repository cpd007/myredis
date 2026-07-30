package core

import (
	"reflect"
	"testing"
)

func Test_readSimpleString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    string
		want2   int
		wantErr bool
	}{
		{
			name:    "test 1",
			data:    []byte("+PONG\r\n"),
			want:    "PONG",
			want2:   7,
			wantErr: false,
		},
		{
			name:    "test 2",
			data:    []byte("+\r\n"),
			want:    "",
			want2:   3,
			wantErr: false,
		},
		{
			name:    "test 3",
			data:    []byte("+PONG"),
			want:    "",
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := readSimpleString(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readSimpleString() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr && gotErr != nil {
				t.Fatal("readSimpleString() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("readSimpleString() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("readSimpleString() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_readError(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    string
		want2   int
		wantErr bool
	}{
		{
			name:    "valid error reply",
			data:    []byte("-ERR unknown command\r\n"),
			want:    "ERR unknown command",
			want2:   22,
			wantErr: false,
		},
		{
			name:    "empty error string",
			data:    []byte("-\r\n"),
			want:    "",
			want2:   3,
			wantErr: false,
		},
		{
			name:    "invalid format",
			data:    []byte("-ERR unknown command"),
			want:    "",
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := readError(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readError() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("readError() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("readError() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("readError() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_readInteger(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    int64
		want2   int
		wantErr bool
	}{
		{
			name:    "zero integer reply",
			data:    []byte(":\r\n"),
			want:    0,
			want2:   3,
			wantErr: false,
		},
		{
			name:    "positive integer reply",
			data:    []byte(":12345\r\n"),
			want:    12345,
			want2:   8,
			wantErr: false,
		},
		{
			name:    "invalid integer format",
			data:    []byte(":12345"),
			want:    0,
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := readInteger(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readInteger() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("readInteger() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("readInteger() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("readInteger() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_readBulkStrings(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    string
		want2   int
		wantErr bool
	}{
		{
			name:    "non-empty bulk string reply",
			data:    []byte("$5\r\nhello\r\n"),
			want:    "hello",
			want2:   11,
			wantErr: false,
		},
		{
			name:    "empty bulk string reply",
			data:    []byte("$0\r\n\r\n"),
			want:    "",
			want2:   6,
			wantErr: false,
		},
		{
			name:    "invalid bulk string length",
			data:    []byte("$x\r\nhello\r\n"),
			want:    "",
			want2:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := readBulkStrings(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readBulkStrings() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("readBulkStrings() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("readBulkStrings() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("readBulkStrings() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_readArray(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    []any
		want2   int
		wantErr bool
	}{
		{
			name:    "array with simple string and integer",
			data:    []byte("*2\r\n+OK\r\n:123\r\n"),
			want:    []any{"OK", int64(123)},
			want2:   15,
			wantErr: false,
		},
		{
			name:    "empty array reply",
			data:    []byte("*0\r\n"),
			want:    []any{},
			want2:   4,
			wantErr: false,
		},
		{
			name:    "invalid array format",
			data:    []byte("*1\r\n+OK"),
			want:    nil,
			want2:   0,
			wantErr: true,
		},
		{
			name:    "nested array",
			data:    []byte("*2\r\n*2\r\n+hello\r\n+world\r\n:42\r\n"),
			want:    []any{[]any{"hello", "world"}, int64(42)},
			want2:   29,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := readArray(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readArray() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("readArray() succeeded unexpectedly")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("readArray() length = %d, want %d", len(got), len(tt.want))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readArray() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("readArray() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
