package cloud

import (
	"time"
)

// HeartBeatReq emmm
type HeartBeatReq struct {
	RaspId        string `json:"rasp_id"`
	PluginVersion string `json:"plugin_version"`
	PluginMd5     string `json:"plugin_md5"`
	ConfigTime    int64  `json:"config_time"`
}

// HeartBeatResp emmm
type HeartBeatResp struct {
	Plugin     *Plugin                 `json:"plugin"`
	Config     *map[string]interface{} `json:"config"`
	ConfigTime int64                   `json:"config_time"`
}

// HeartBeat emmm
func (c *Client) HeartBeat(updatePlugin func(string, string), updateConfig func(*map[string]interface{})) error {
	request := HeartBeatReq{
		c.rasp.Id,
		c.plugin.Version,
		c.plugin.Md5,
		c.configTime,
	}
	var response HeartBeatResp
	if err := c.Post("/v1/agent/heartbeat", &request, &response); err != nil {
		return err
	}
	if response.Plugin != nil && c.plugin.Md5 != response.Plugin.Md5 {
		c.plugin = *response.Plugin
		updatePlugin(c.plugin.Content, c.plugin.Name)
	}
	if response.Config != nil {
		c.config = *response.Config
		c.configTime = response.ConfigTime
		updateConfig(&c.config)
	}
	return nil
}

// StartHeartBeat emmm
func (c *Client) StartHeartBeat(interval time.Duration, updatePlugin func(string, string), updateConfig func(*map[string]interface{}), onError func(error)) {
	c.wg.Add(1)
	c.isHeartBeat = true
ABORT:
	for {
		select {
		case <-time.After(interval):
			err := c.HeartBeat(updatePlugin, updateConfig)
			if err != nil {
				onError(err)
				time.Sleep(1 * time.Second)
				c.HeartBeat(updatePlugin, updateConfig)
			}
		case <-c.abort:
			break ABORT
		}
	}
	c.wg.Done()
	c.isHeartBeat = false
}

// StopHeartBeat emmm
func (c *Client) StopHeartBeat() {
	if c.isHeartBeat {
		c.abort <- struct{}{}
		c.wg.Wait()
	}
}
