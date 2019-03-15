package config

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type UpdateListener interface {
	OnConfigUpdate()
}

type GeneralConfig struct {
	general   *viper.Viper
	listeners []UpdateListener
	mu        sync.RWMutex
}

func NewGeneralConfig() *GeneralConfig {
	generalViper := viper.New()
	generalViper.SetDefault("plugin.timeout.millis", 100)
	generalViper.SetDefault("plugin.maxstack", 100)
	generalViper.SetDefault("plugin.filter", false)
	generalViper.SetDefault("log.maxburst", 100)
	generalViper.SetDefault("log.maxstack", 10)
	generalViper.SetDefault("log.maxbackup", 30)
	generalViper.SetDefault("syslog.tag", "OpenRASP")
	generalViper.SetDefault("syslog.url", "")
	generalViper.SetDefault("syslog.facility", 1)
	generalViper.SetDefault("syslog.enable", false)
	generalViper.SetDefault("syslog.connection_timeout", 50)
	generalViper.SetDefault("syslog.read_timeout", 10)
	generalViper.SetDefault("syslog.reconnect_interval", 300)
	generalViper.SetDefault("block.status_code", 302)
	generalViper.SetDefault("block.redirect_url", `https://rasp.baidu.com/blocked/?request_id=%request_id%`)
	generalViper.SetDefault("block.content_json", `{"error":true, "reason": "Request blocked by OpenRASP", "request_id": "%request_id%"}`)
	generalViper.SetDefault("block.content_xml", `<?xml version="1.0"?><doc><error>true</error><reason>Request blocked by OpenRASP</reason><request_id>%request_id%</request_id></doc>`)
	generalViper.SetDefault("block.content_html", `</script><script>location.href="https://rasp.baidu.com/blocked2/?request_id=%request_id%"</script>`)
	generalViper.SetDefault("inject.urlprefix", "")
	generalViper.SetDefault("inject.custom_headers", []string{})
	generalViper.SetDefault("body.maxbytes", 4096)
	generalViper.SetDefault("clientip.header", "")
	generalViper.SetDefault("security.enforce_policy", false)
	generalViper.SetDefault("lru.max_size", 1024)
	generalViper.SetDefault("hook.white", map[string]interface{}{})
	generalViper.SetDefault("decompile.enable", false)
	return &GeneralConfig{
		general: generalViper,
	}
}

func (gc *GeneralConfig) AttachListener(listener UpdateListener) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	gc.listeners = append(gc.listeners, listener)
}

func (gc *GeneralConfig) GetBool(key string) bool {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.general.GetBool(key)
}

func (gc *GeneralConfig) GetString(key string) string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.general.GetString(key)
}

func (gc *GeneralConfig) GetInt(key string) int {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.general.GetInt(key)
}

func (gc *GeneralConfig) GetInt64(key string) int64 {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.general.GetInt64(key)
}

func (gc *GeneralConfig) ReadConfig(in io.Reader) (err error) {
	gc.mu.Lock()
	defer func() {
		gc.mu.Unlock()
		if err == nil {
			for _, l := range gc.listeners {
				l.OnConfigUpdate()
			}
		}
	}()
	gc.general.SetConfigType("yaml")
	err = gc.general.ReadConfig(in)
	if err != nil {
		log.Printf("%v", err)
	}
	return err
}

func (gc *GeneralConfig) OnUpdate(absPath string) {
	raw, err := os.Open(absPath)
	if err == nil {
		gc.ReadConfig(raw)
	}
}
