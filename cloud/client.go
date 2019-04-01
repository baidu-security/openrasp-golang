package cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Client emmm
type Client struct {
	http.Client
	abort       chan struct{}
	isHeartBeat bool
	host        string
	appid       string
	appsecret   string
	rasp        Rasp
	plugin      Plugin
	config      map[string]interface{}
	configTime  int64
}

// NewClient emmm
func NewClient(host, appid, appsecret string, timeout time.Duration) *Client {
	c := &Client{
		Client:    http.Client{Timeout: timeout},
		abort:     make(chan struct{}),
		host:      host,
		appid:     appid,
		appsecret: appsecret,
	}
	return c
}

// Post emmm
func (c *Client) Post(path string, request, response interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.host+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OpenRASP-AppID", c.appid)
	req.Header.Set("X-OpenRASP-AppSecret", c.appsecret)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	var temp struct {
		Status      int         `json:"status"`
		Description string      `json:"description"`
		Data        interface{} `json:"data"`
	}
	temp.Data = response
	if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
		return err
	}
	if temp.Status != 0 {
		return errors.New(temp.Description)
	}
	return nil
}
