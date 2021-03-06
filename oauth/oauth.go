package oauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/x-io/wechat/config"
	"github.com/x-io/wechat/param"
	"github.com/x-io/wechat/util"
)

//OAuth OAuth
type OAuth struct {
}

//New OAuth
func New() *OAuth {
	return &OAuth{}
}

//GetMiniSession 通过小程序授权的code 换取session_key
func (o *OAuth) GetMiniSession(key, code string) (param.Params, error) {
	config, err := config.GetConfig(key)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(jscode2SessionURL, config.AppID, config.AppSecret, code)
	params := make(param.Params)
	if err := getBind(url, params); err != nil {
		return nil, err
	}

	return params, nil
}

//GetMiniDecrypt 小程序数据解密 通过 session_key 解密数据
//
// 参数: session: 用户 session_key; iv: iv数据; data: 解密数据
//
// 返回: param, error
func (o *OAuth) GetMiniDecrypt(session, iv, data string) (param.Params, error) {
	_key, err := base64.StdEncoding.DecodeString(session)
	if err != nil {
		return nil, err
	}
	_iv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	_data, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	_data, err = util.AesCBCDecrypt2(_data, _key, _iv)
	if err != nil {
		return nil, err
	}

	params := make(param.Params)
	if err = json.Unmarshal(_data, &params); err != nil {
		return nil, err
	}

	return params, nil
}

//GetMiniDecryptByCode 小程序数据解密 通过 code 解密数据
//
// 参数: key: 微信账号; code: 小程序端获取到的code; iv: iv数据; data: 解密数据
//
// 返回: param, error
func (o *OAuth) GetMiniDecryptByCode(key, code, iv, data string) (param.Params, error) {
	result, err := o.GetMiniSession(key, code)
	if err != nil {
		return nil, err
	}

	_data, err := o.GetMiniDecrypt(result.Get("session_key"), iv, data)
	if err != nil {
		return nil, err
	}
	_data["openid"] = result.Get("openid")
	return _data, nil
}

//GetURL 获取跳转的url地址
func (o *OAuth) GetURL(key, redirectURI, scope, state string) (string, error) {
	config, err := config.GetConfig(key)
	if err != nil {
		return "", err
	}

	//url encode
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(redirectOauthURL, config.AppID, urlStr, scope, state), nil
}

//GetToken 通过网页授权的code 换取access_token
func (o *OAuth) GetToken(key, code string) (param.Params, error) {
	config, err := config.GetConfig(key)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(accessTokenURL, config.AppID, config.AppSecret, code)
	params := make(param.Params)
	if err := getBind(url, params); err != nil {
		return nil, err
	}
	//fmt.Println(params)
	return params, nil
}

//RefreshOAutnToken 刷新access_token
func (o *OAuth) RefreshOAutnToken(key, refreshToken string) (param.Params, error) {
	config, err := config.GetConfig(key)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(refreshAccessTokenURL, config.AppID, refreshToken)

	params := make(param.Params)

	if err := getBind(url, params); err != nil {
		return nil, err
	}

	return params, nil
}

//CheckOAutnToken 检验access_token是否有效
func (o *OAuth) CheckOAutnToken(accessToken, openID string) (b bool, err error) {

	url := fmt.Sprintf(checkAccessTokenURL, accessToken, openID)
	if _, err := util.HTTPGet(url); err != nil {
		return false, err
	}

	return true, nil
}

//GetUserData 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (o *OAuth) GetUserData(accessToken, openID string) (response []byte, err error) {
	url := fmt.Sprintf(userInfoURL, accessToken, openID)

	return util.HTTPGet(url)
}

//GetUserInfo 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (o *OAuth) GetUserInfo(accessToken, openID string) (param.Params, []byte, error) {

	url := fmt.Sprintf(userInfoURL, accessToken, openID)
	var data param.Params

	response, err := getJSONBind(url, &data)

	if err != nil {
		return nil, nil, err
	}

	return data, response, nil
}

//GetInfo GetInfo
func (o *OAuth) GetInfo(key, code string) (string, []byte, error) {

	result, err := o.GetToken(key, code)
	if err != nil {
		return "", nil, err
	}

	openID := result.Get("openid")
	// if result.Get("scope") == "snsapi_base" {
	// 	return openID, nil, nil
	// }

	info, err := o.GetUserData(result.Get("access_token"), openID)
	if err != nil {
		return openID, nil, nil
	}
	return openID, info, err
}

func getBind(url string, v param.Params) error {
	response, err := util.HTTPGet(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(response, &v)
	if err != nil {
		return err
	}

	return nil
}

func getJSONBind(url string, v interface{}) (response []byte, err error) {
	response, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, v)
	if err != nil {
		return
	}
	return
}
