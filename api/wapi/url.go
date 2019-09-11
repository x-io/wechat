package wapi

import (
	"encoding/json"

	"github.com/x-io/wechat/util"
)

//WeURL WeURL
type WeURL string

const (
	//AccessTokenURL 获取access_token的接口地址
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	//TicketURL 获取jsapi_ticket的接口地址
	TicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	//TemplateSendURL 设置模板的接口地址
	TemplateSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

func getBind(url string, v interface{}) error {
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
	// if result.ErrCode != 0 {
	// 	err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
	// 	return
	// }

	return
}
