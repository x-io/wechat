package api

import (
	"fmt"

	"github.com/x-io/wechat/api/wapi"
	"github.com/x-io/wechat/cache"
)

//GetAccessToken GetAccessToken
func (m *API) GetAccessToken(key string) (accessToken string, err error) {
	cacheKey := fmt.Sprintf("access_token_%s", key)
	token, err := cache.Get(cacheKey)
	if err == nil {
		return token, nil
	}

	m.accessTokenLock.Lock()
	defer m.accessTokenLock.Unlock()

	token, err = cache.Get(cacheKey)
	if err == nil {
		return token, nil
	}

	config, err := cache.GetConfig(key)
	if err != nil {
		return "", err
	}

	//从微信服务器获取
	data, err := wapi.GetAccessToken(config.AppID, config.AppSecret)
	if err != nil {
		return
	}
	//fmt.Println("AccessToken", data.AccessToken, data.ExpiresIn)
	cache.Set(cacheKey, data.AccessToken, data.ExpiresIn)

	accessToken = data.AccessToken
	return
}
