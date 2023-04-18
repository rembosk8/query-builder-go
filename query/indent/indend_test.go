package indent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const sFormat = "#%s#"

type mockSanitizer struct {
}

func (m mockSanitizer) Sanitize(v string) string {
	return fmt.Sprintf(sFormat, v)
}

var _ Sanitizer = mockSanitizer{}

func TestBuilder_Build(t *testing.T) {
	sanitizer := new(mockSanitizer)
	bdr := NewBuilder(sanitizer)

	name := "table"

	assert.Equal(t, fmt.Sprintf(sFormat, name), bdr.Indent(name).String())

	bdr = NewBuilder(nil)
	assert.Equal(t, name, bdr.Indent(name).String())
}
