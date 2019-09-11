package wapi

import (
	"fmt"
)

//AccessToken AccessToken
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int  `json:"expires_in"`
}

//GetAccessToken 强制从微信服务器获取token
func GetAccessToken(appID, appSecret string) (*AccessToken, error) {
	url := fmt.Sprintf(AccessTokenURL, appID, appSecret)
	var data AccessToken

	if err := getBind(url, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
