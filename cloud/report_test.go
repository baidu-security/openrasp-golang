package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	c := NewClient("http://scloud.baidu.com:8087", "043b6f2ad443a858ca1b3c593b9baa12e75b6041", "jVRYpHVfkBq6XF0Wb73hX1AvN3Xo9ZiSxKY6xdBQVee", 20*time.Second)
	c.Register("569e8ea7a16123492b5878920fd36985", "/tmp", "tmp", "golang", "1", 60)
	err := c.Report(1)
	assert.NoError(t, err)
}
