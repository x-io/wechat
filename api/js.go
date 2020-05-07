package api

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/x-io/cache"
	"github.com/x-io/wechat/api/wapi"
	"github.com/x-io/wechat/config"
	"github.com/x-io/wechat/util"
)

// JSConfig 返回给用户jssdk配置信息
type JSConfig struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

//GetJSConfig GetJSConfig
func (m *API) GetJSConfig(key, uri string) (jsconfig *JSConfig, err error) {

	var ticketStr string
	ticketStr, err = m.GetTicket(key)

	if err != nil {
		return nil, err
	}

	config, err := config.GetConfig(key)
	if err != nil {
		return nil, err
	}

	nonceStr := util.RandomStr(16)
	timestamp := util.GetCurrTs()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	sigStr := jsSignature(str)
	jsconfig = new(JSConfig)
	jsconfig.AppID = config.AppID
	jsconfig.NonceStr = nonceStr
	jsconfig.Timestamp = timestamp
	jsconfig.Signature = sigStr
	return
}

//GetTicket 获取jsapi_ticket
func (m *API) GetTicket(key string) (ticketStr string, err error) {

	cacheKey := fmt.Sprintf("jsapi_ticket_%s", key)
	token, err := cache.Get(cacheKey)

	if err == nil {
		return string(token), nil
	}

	m.ticketLock.Lock()
	defer m.ticketLock.Unlock()
	token, err = cache.Get(cacheKey)
	if err == nil {
		return string(token), nil
	}

	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return "", err
	}

	//从微信服务器获取
	data, err := wapi.GetTicket(accessToken)
	if err != nil {
		return "", err
	}
	ticketStr = data.Ticket
	cache.Set(cacheKey, []byte(data.Ticket), time.Duration(data.ExpiresIn)*time.Second)

	return
}

//Signature sha1签名
func jsSignature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
