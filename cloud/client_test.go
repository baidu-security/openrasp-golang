package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient("http://scloud.baidu.com:8087", "", "", 20*time.Second)
	assert.NotEmpty(t, c)
}

func TestPost(t *testing.T) {
	c := NewClient("http://scloud.baidu.com:8087", "", "", 20*time.Second)
	var resp interface{}
	err := c.Post("/v1/agent/rasp", &struct{}{}, &resp)
	assert.Error(t, err)
	assert.EqualError(t, err, "Unauthorized")
}
