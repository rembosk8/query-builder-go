package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getNums(t *testing.T) {
	type args struct {
		s   int
		cnt int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1-39",
			args: args{
				s:   1,
				cnt: 39,
			},
			want: "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39",
		},
		{
			name: "2-3",
			args: args{
				s:   2,
				cnt: 3,
			},
			want: "$2, $3, $4",
		},
		{
			name: "10-13",
			args: args{
				s:   10,
				cnt: 4,
			},
			want: "$10, $11, $12, $13",
		},
		{
			name: "11-15",
			args: args{
				s:   11,
				cnt: 5,
			},
			want: "$11, $12, $13, $14, $15",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getNums(tt.args.s, tt.args.cnt), "getNums(%v, %v)", tt.args.s, tt.args.cnt)
		})
	}
}
