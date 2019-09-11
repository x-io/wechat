package api

import (
	"encoding/json"

	"github.com/x-io/wechat/api/wapi"
)

//Message 发送的模板消息内容
type Message struct {
	ToUser     string              `json:"touser"`          // 必须, 接受者OpenID
	TemplateID string              `json:"template_id"`     // 必须, 模版ID
	URL        string              `json:"url,omitempty"`   // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	Color      string              `json:"color,omitempty"` // 可选, 整个消息的颜色, 可以不设置
	Data       map[string]DataItem `json:"data"`            // 必须, 模板数据

	MiniProgram struct {
		AppID    string `json:"appid"`    //所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
		PagePath string `json:"pagepath"` //所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar）
	} `json:"miniprogram,omitempty"` //可选,跳转至小程序地址
}

//DataItem 模版内某个 .DATA 的值
type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

//SendTemplate 获取用户基本信息
func (m *API) SendTemplate(key string, message *Message) (msgID int64, err error) {

	var accessToken string
	accessToken, err = m.GetAccessToken(key)
	if err != nil {
		return
	}

	var data []byte
	data, err = json.Marshal(message)
	if err != nil {
		return
	}

	return wapi.SendTemplate(accessToken, data)
}
