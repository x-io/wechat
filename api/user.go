package api

import (
	"encoding/json"

	"github.com/x-io/wechat/api/wapi"
)

//Info 用户基本信息
type Info struct {
	Subscribe     int32   `json:"subscribe"`
	OpenID        string  `json:"openid"`
	Nickname      string  `json:"nickname"`
	Sex           int32   `json:"sex"`
	City          string  `json:"city"`
	Country       string  `json:"country"`
	Province      string  `json:"province"`
	Language      string  `json:"language"`
	Headimgurl    string  `json:"headimgurl"`
	SubscribeTime int32   `json:"subscribe_time"`
	UnionID       string  `json:"unionid"`
	Remark        string  `json:"remark"`
	GroupID       int32   `json:"groupid"`
	TagidList     []int32 `json:"tagid_list"`
}

//GetUserInfo 获取用户基本信息
func (m *API) GetUserInfo(key, openID string) (userInfo *Info, err error) {
	var accessToken string
	accessToken, err = m.GetAccessToken(key)

	if err != nil {

		return
	}

	var response []byte
	response, err = wapi.GetUserInfo(accessToken, openID)
	if err != nil {
		return
	}
	//fmt.Println(string(response))
	userInfo = new(Info)
	err = json.Unmarshal(response, userInfo)
	if err != nil {
		return
	}

	return
}

//GetUserJSON 获取用户原始信息
func (m *API) GetUserJSON(key, openID string) ([]byte, error) {

	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.GetUserInfo(accessToken, openID)
}

// // UpdateRemark 设置用户备注名
// func (m *API) UpdateRemark(openID, remark string) (err error) {
// 	var accessToken string
// 	accessToken, err = m.GetAccessToken()
// 	if err != nil {
// 		return
// 	}

// 	uri := fmt.Sprintf(updateRemarkURL, accessToken)
// 	var response []byte
// 	response, err = util.PostJSON(uri, map[string]string{"openid": openID, "remark": remark})
// 	if err != nil {
// 		return
// 	}

// 	return util.DecodeWithCommonError(response, "UpdateRemark")
// }
