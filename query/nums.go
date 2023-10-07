package query

import "fmt"

const (
	oneDigNums = 9
	numCnt     = 59
	numsStr    = " $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59" //nolint:lll
)

func genNums(s, cnt int) string {
	res := fmt.Sprintf("$%d", s)
	for i := 1; i < cnt; i++ {
		res += fmt.Sprintf(", $%d", s+i)
	}

	return res
}

func getNums(s, cnt int) string {
	if s == 0 {
		return ""
	}
	if s+cnt-1 > numCnt {
		return genNums(s, cnt)
	}
	s--

	start := s * 4               //nolint:gomnd
	pos := start + (cnt-1)*4 + 3 //nolint:gomnd

	if s > oneDigNums {
		start += s - oneDigNums
	}
	if s+cnt > oneDigNums {
		pos += s + cnt - oneDigNums
	}

	return numsStr[start+1 : pos]
}
