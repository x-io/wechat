package wapi

import (
	"fmt"

	"github.com/x-io/wechat/util"
)

const (
	qrcodeAURL = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	qrcodeBURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
	qrcodeCURL = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=%s"
)

//GetQRCodeByA 获取小程序二维码
func GetQRCodeByA(accessToken string, body []byte) (response []byte, err error) {
	uri := fmt.Sprintf(qrcodeAURL, accessToken)
	return util.HTTPPost(uri, body)
}

//GetQRCodeByB 获取小程序二维码
func GetQRCodeByB(accessToken string, body []byte) (response []byte, err error) {
	uri := fmt.Sprintf(qrcodeBURL, accessToken)
	return util.HTTPPost(uri, body)
}

//GetQRCodeByC 获取小程序二维码
func GetQRCodeByC(accessToken string, body []byte) (response []byte, err error) {
	uri := fmt.Sprintf(qrcodeCURL, accessToken)
	return util.HTTPPost(uri, body)
}
