package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeartBeat(t *testing.T) {
	c := NewClient("http://scloud.baidu.com:8087", "043b6f2ad443a858ca1b3c593b9baa12e75b6041", "jVRYpHVfkBq6XF0Wb73hX1AvN3Xo9ZiSxKY6xdBQVee", 20*time.Second)
	c.Register("569e8ea7a16123492b5878920fd36985", "/tmp", "tmp", "golang", "1", 1)
	err := c.HeartBeat(func(source, filename string) {
		assert.NotEmpty(t, source)
		assert.NotEmpty(t, filename)
	}, func(config *map[string]interface{}) {
		assert.NotEmpty(t, config)
	})
	assert.NoError(t, err)
	go c.StartHeartBeat(1*time.Second, func(source, filename string) {
		t.Fatal()
	}, func(config *map[string]interface{}) {
		t.Fatal()
	})
	time.Sleep(5 * time.Second)
	assert.True(t, c.isHeartBeat)
	c.StopHeartBeat()
	assert.False(t, c.isHeartBeat)
}
