package sanitize

import (
	"testing"
	"time"
)

func TestValue_Sanitize(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{
			name: "nil",
			val:  nil,
			want: "null",
		},
		{
			name: "int64",
			val:  int64(100),
			want: "100",
		},
		{
			name: "fload64",
			val:  100.99999,
			want: "100.99999",
		},
		{
			name: "bool",
			val:  true,
			want: "true",
		},
		{
			name: "bytes",
			val:  []byte{1, 2, 3, 0xff},
			want: "'\\x010203ff'",
		},
		{
			name: "string",
			val:  "test",
			want: "'test'",
		},
		{
			name: "time",
			val:  time.Date(2000, 1, 1, 1, 1, 1, 999999, time.UTC),
			want: "'2000-01-01 01:01:01.000999Z'",
		},
		{
			name: "int",
			val:  32,
			want: "32",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Value{}
			if got := v.Sanitize(tt.val); got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
