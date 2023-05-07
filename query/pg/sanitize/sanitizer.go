package sanitize

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const defaultVal = "DEFAULT" //todo: make enum value from that

type Indent struct{}

func (i Indent) Sanitize(val string) string {
	s := strings.ReplaceAll(val, string([]byte{0}), "")
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

type Value struct{}

func (v Value) Sanitize(val any) string {
	var str string
	switch x := val.(type) {
	case nil:
		str = "null"
	case int64:
		str = strconv.FormatInt(x, 10)
	case float64:
		str = strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		str = strconv.FormatBool(x)
	case []byte:
		str = QuoteBytes(x)
	case string:
		if x == defaultVal {
			str = x
		} else {
			str = QuoteString(x)
		}
	case time.Time:
		str = x.Truncate(time.Microsecond).Format("'2006-01-02 15:04:05.999999999Z07:00:00'")
	default:
		str = fmt.Sprintf("%v", x)
	}

	return str
}

func (v Value) IsDefault(val any) bool {
	if v, ok := val.(string); ok {
		return v == defaultVal
	}

	return false
}

func QuoteString(str string) string {
	return "'" + strings.ReplaceAll(str, "'", "''") + "'"
}

func QuoteBytes(buf []byte) string {
	return `'\x` + hex.EncodeToString(buf) + "'"
}
