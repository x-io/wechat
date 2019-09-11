package wapi

import (
	"fmt"

	"github.com/x-io/wechat/util"
)

const (
	userInfoURL     = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
	updateRemarkURL = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=%s"
	userCodeURL     = "https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code"
)

//GetUserInfo 从服务器中获取UserInfo
func GetUserInfo(accessToken, openID string) ([]byte, error) {
	url := fmt.Sprintf(userInfoURL, accessToken, openID)

	return util.HTTPGet(url)
}

// UpdateRemark 设置用户备注名
func UpdateRemark(accessToken, openID, remark string) (err error) {

	uri := fmt.Sprintf(updateRemarkURL, accessToken)

	_, err = util.PostJSON(uri, map[string]string{"openid": openID, "remark": remark})

	return
}
