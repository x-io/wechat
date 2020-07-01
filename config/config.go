package config

import (
	"errors"
	"net/http"
	"sync"
)

//Config Config
type Config struct {
	AppID     string
	AppSecret string
	Pay       struct {
		MchID      string
		APIKey     string
		ServiceID  string
		SignType   string
		CertFile   string
		CertData   string
		SandboxKey string
		Transport  *http.Transport
	}
}

var (
	db map[string]*Config

	query func(key string) (*Config, error)
	lock  sync.RWMutex
)

func init() {
	db = make(map[string]*Config)
}

//SetConfig SetConfig
func SetConfig(call func(key string) (*Config, error)) {
	query = call
}

//GetConfig GetConfig
func GetConfig(key string) (*Config, error) {
	if v, ok := db[key]; ok {
		return v, nil
	}

	lock.Lock()
	defer lock.Unlock()
	if v, ok := db[key]; ok {
		return v, nil
	}

	if query == nil {
		return nil, errors.New("该账户数据不存在")
	}

	config, err := query(key)
	if err != nil {
		return nil, errors.New("该账户数据不存在")
	}

	db[key] = config

	return config, nil
}
