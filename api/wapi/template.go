package wapi

import (
	"encoding/json"
	"fmt"

	"github.com/x-io/wechat/util"
)

const (
	templateSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

//SendTemplate SendTemplate
func SendTemplate(accessToken string, data []byte) (msgID int64, err error) {

	uri := fmt.Sprintf(templateSendURL, accessToken)
	var response []byte
	response, err = util.HTTPPost(uri, data)
	if err != nil {
		return
	}

	var result struct {
		MsgID int64 `json:"msgid"`
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}

	msgID = result.MsgID
	return
}
