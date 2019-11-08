package cache

import (
	"errors"
	"net/http"
	"sync"

	"github.com/coocood/freecache"
)

//Config Config
type Config struct {
	AppID     string
	AppSecret string
	Pay       struct {
		MchID    string
		APIKey   string
		SignType string
		Cert     []byte

		SandboxKey string
		Transport  *http.Transport
	}
}

var (
	db    map[string]*Config
	cache *freecache.Cache

	query func(key string) (*Config, error)
	lock  sync.RWMutex
)

func init() {

	db = make(map[string]*Config)
	cache = freecache.NewCache(1024)
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

//Set Cache Set
func Set(key, value string, expireSeconds int) error {
	return cache.Set([]byte(key), []byte(value), expireSeconds)
}

//Get Cache Get
func Get(key string) (string, error) {
	var val []byte
	val, err := cache.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(val), nil
}
