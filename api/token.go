package api

import (
	"fmt"
	"time"

	"github.com/x-io/cache"
	"github.com/x-io/wechat/api/wapi"
	"github.com/x-io/wechat/config"
)

//GetAccessToken GetAccessToken
func (m *API) GetAccessToken(key string) (accessToken string, err error) {
	cacheKey := fmt.Sprintf("access_token_%s", key)
	token, err := cache.Get(cacheKey)
	if err == nil {
		return string(token), nil
	}

	m.accessTokenLock.Lock()
	defer m.accessTokenLock.Unlock()

	token, err = cache.Get(cacheKey)
	if err == nil {
		return string(token), nil
	}

	config, err := config.GetConfig(key)
	if err != nil {
		return "", err
	}

	//从微信服务器获取
	data, err := wapi.GetAccessToken(config.AppID, config.AppSecret)
	if err != nil {
		return
	}
	//fmt.Println("AccessToken", data.AccessToken, data.ExpiresIn)
	cache.Set(cacheKey, []byte(data.AccessToken), time.Duration(data.ExpiresIn)*time.Second)

	accessToken = data.AccessToken
	return
}
