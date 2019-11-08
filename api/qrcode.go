package api

import (
	"encoding/json"

	"github.com/x-io/wechat/api/wapi"
)

//GetQRCodeByA 获取A型小程序二维码
// 通过该接口生成的小程序码，永久有效，有数量限制
func (m *API) GetQRCodeByA(key string, config map[string]interface{}) ([]byte, error) {
	var err error
	var accessToken string
	if accessToken, err = m.GetAccessToken(key); err != nil {
		return nil, err
	}

	var body []byte
	if body, err = json.Marshal(config); err != nil {
		return nil, err
	}

	return wapi.GetQRCodeByA(accessToken, body)
}

//GetQRCodeByB 获取B型小程序二维码
// 通过该接口生成的小程序码，永久有效，数量暂无限制
func (m *API) GetQRCodeByB(key string, config map[string]interface{}) ([]byte, error) {
	var err error
	var accessToken string
	if accessToken, err = m.GetAccessToken(key); err != nil {
		return nil, err
	}

	var body []byte
	if body, err = json.Marshal(config); err != nil {
		return nil, err
	}

	return wapi.GetQRCodeByB(accessToken, body)
}

//GetQRCodeByC 获取C型小程序二维码
// 通过该接口生成的小程序码，永久有效，有数量限制
func (m *API) GetQRCodeByC(key string, config map[string]interface{}) ([]byte, error) {
	var err error
	var accessToken string
	if accessToken, err = m.GetAccessToken(key); err != nil {
		return nil, err
	}

	var body []byte
	if body, err = json.Marshal(config); err != nil {
		return nil, err
	}

	return wapi.GetQRCodeByC(accessToken, body)
}
