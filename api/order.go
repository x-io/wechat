package api

import "github.com/x-io/wechat/api/wapi"

//GetOrderByID 根据订单号获取小店订单信息
func (m *API) GetOrderByID(key, ordrID string) ([]byte, error) {

	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.GetOrderByID(accessToken, ordrID)
}

//GetOrderList 按条件获取订单列表
func (m *API) GetOrderList(key string, status, beginTime, endTime int) ([]byte, error) {
	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.GetOrderList(accessToken, status, beginTime, endTime)
}

//SetDelivery 将订单发货处理
func (m *API) SetDelivery(key, ordrID string) ([]byte, error) {

	accessToken, err := m.GetAccessToken(key)
	if err != nil {
		return nil, err
	}

	return wapi.SetDelivery(accessToken, ordrID)
}
