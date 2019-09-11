package api

import (
	"fmt"
	"sync"

	"github.com/x-io/wechat/api/wapi"
	"github.com/x-io/wechat/cache"
)

//API API Manager
type API struct {
	accessTokenLock sync.RWMutex
	ticketLock      sync.RWMutex
}

//New API
func New() *API {
	return &API{}
}

//GetAccessToken GetAccessToken
func (m *API) GetAccessToken(key string) (accessToken string, err error) {
	//return "16_qzabqubXMdYNeRBt3xA4v4YkcxiyDFiED1VzGZghfrpb8IXblR0lV_liPjJ5SN2Q2wPl2GqQvB97lWLnSdoLeNEcWph7Y49P_WL0b0EigT66wfcdYPKhjUjIPJP4ODPfUeNwubPqGU7-ZtkSCQUiADAWKO", nil
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
	fmt.Println("AccessToken", data.AccessToken, data.ExpiresIn)
	cache.Set(cacheKey, data.AccessToken, data.ExpiresIn)

	accessToken = data.AccessToken
	return
}
