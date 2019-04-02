package cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client emmm
type Client struct {
	http.Client
	abort       chan struct{}
	isHeartBeat bool
	wg          sync.WaitGroup
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
	body, err := c.PostRaw(path, data)
	if err != nil {
		return err
	}
	defer body.Close()
	var temp struct {
		Status      int         `json:"status"`
		Description string      `json:"description"`
		Data        interface{} `json:"data"`
	}
	temp.Data = response
	if err := json.NewDecoder(body).Decode(&temp); err != nil {
		return err
	}
	if temp.Status != 0 {
		return errors.New(temp.Description)
	}
	return nil
}

// PostRaw emmm
func (c *Client) PostRaw(path string, request []byte) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", c.host+path, bytes.NewReader(request))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OpenRASP-AppID", c.appid)
	req.Header.Set("X-OpenRASP-AppSecret", c.appsecret)
	resp, err := c.Do(req)
	return resp.Body, err
}
