package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	c := NewClient("http://scloud.baidu.com:8087", "043b6f2ad443a858ca1b3c593b9baa12e75b6041", "jVRYpHVfkBq6XF0Wb73hX1AvN3Xo9ZiSxKY6xdBQVee", 20*time.Second)
	c.Register("569e8ea7a16123492b5878920fd36985", "/tmp", "tmp", "golang", "1", 60)
	log := map[string]interface{}{
		"attack_type": "directory",
	}
	logs := []map[string]interface{}{log}
	count, err := c.Log("attack", logs)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
