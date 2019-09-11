package wapi

import (
	"bytes"
	"fmt"

	"github.com/x-io/wechat/util"
)

const (
	orderByIDURL     = "https://api.weixin.qq.com/merchant/order/getbyid?access_token=%s"
	orderListURL     = "https://api.weixin.qq.com/merchant/order/getbyfilter?access_token=%s"
	orderDeliveryURL = "https://api.weixin.qq.com/merchant/order/setdelivery?access_token=%s"
)

//GetOrderByID 根据订单号获取订单
func GetOrderByID(accessToken, ordrID string) (response []byte, err error) {

	uri := fmt.Sprintf(orderByIDURL, accessToken)
	body := fmt.Sprintf(`{"order_id": "%s"}`, ordrID)

	response, err = util.HTTPPost(uri, []byte(body))
	return
}

//GetOrderList 获取订单列表
func GetOrderList(accessToken string, status, beginTime, endTime int) (response []byte, err error) {

	uri := fmt.Sprintf(orderListURL, accessToken)
	var buffer bytes.Buffer
	buffer.WriteString("{")
	if status > -1 {
		buffer.WriteString(`"status":` + string(status) + `,`)
	}

	if beginTime > 0 {
		buffer.WriteString(`"begintime":` + string(beginTime) + `,`)
	}

	if endTime > 0 {
		buffer.WriteString(`"endtime":` + string(endTime) + `,`)
	}
	buffer.WriteString("}")

	response, err = util.HTTPPost(uri, buffer.Bytes())
	return
}

//SetDelivery 订单发货
func SetDelivery(accessToken, ordrID string) (response []byte, err error) {
	uri := fmt.Sprintf(orderDeliveryURL, accessToken)
	body := fmt.Sprintf(`{"order_id": "%s","need_delivery": 0}`, ordrID)

	response, err = util.HTTPPost(uri, []byte(body))
	return
}
