package pay

import (
	"errors"
	"strconv"

	"github.com/x-io/wechat/cache"
	"github.com/x-io/wechat/param"
	"github.com/x-io/wechat/util"
)

//Params Params
type Params = param.Params

//Client Pay Client
type Client struct {
	isSandbox bool
}

//New 创建微信支付客户端
func New() *Client {
	return &Client{}
}

//Sandbox 开启沙箱模式
func (c *Client) Sandbox() error {
	c.isSandbox = true

	return nil
}

//ValidNotify 交易通知验证
func (c *Client) ValidNotify(key string, body string) (Params, error) {

	params := Params(util.XMLToMap(body))

	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	if validSign(config.Pay.SignType, config.Pay.APIKey, params) {
		return params, nil
	}

	return params, errors.New("wechat: 签名认证失败")
}

//ChooseWXPay 获取Js调用参数
func (c *Client) ChooseWXPay(key string, prepayID string) (Params, error) {
	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	params := make(Params).Set("timeStamp", strconv.FormatInt(util.GetCurrTs(), 10))
	params["appId"] = config.AppID
	params["signType"] = config.Pay.SignType
	params["nonceStr"] = util.NonceStr()
	params["package"] = "prepay_id=" + prepayID
	params["paySign"] = sign(config.Pay.SignType, config.Pay.APIKey, params)
	return params, nil
}

//UnifiedOrder 统一下单
func (c *Client) UnifiedOrder(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = unifiedOrderSandboxURL
	} else {
		url = unifiedOrderURL
	}
	return sendAPI(key, url, params, false)
}

//MicroPay 刷卡支付
func (c *Client) MicroPay(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = microPaySandboxURL
	} else {
		url = microPayURL
	}
	return sendAPI(key, url, params, false)
}

//Refund 退款
func (c *Client) Refund(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = refundSandboxURL
	} else {
		url = refundURL
	}
	return sendAPI(key, url, params, true)
}

//Reverse 撤销订单
func (c *Client) Reverse(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = reverseSandboxURL
	} else {
		url = reverseURL
	}
	return sendAPI(key, url, params, true)
}

//CloseOrder 关闭订单
func (c *Client) CloseOrder(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = closeOrderSandboxURL
	} else {
		url = closeOrderURL
	}
	return sendAPI(key, url, params, false)
}

//OrderQuery 订单查询
func (c *Client) OrderQuery(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = orderQuerySandboxURL
	} else {
		url = orderQueryURL
	}
	return sendAPI(key, url, params, false)
}

//RefundQuery 退款查询
func (c *Client) RefundQuery(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = refundQuerySandboxURL
	} else {
		url = refundQueryURL
	}
	return sendAPI(key, url, params, false)
}

//DownloadBill 对账单下载
func (c *Client) DownloadBill(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = downloadBillSandboxURL
	} else {
		url = downloadBillURL
	}
	return sendAPI(key, url, params, false)

	// var p Params

	// // 如果出现错误，返回XML数据
	// if strings.Index(xmlStr, "<") == 0 {
	// 	p = XmlToMap(xmlStr)
	// 	return p, err
	// } else { // 正常返回csv数据
	// 	p.SetString("return_code", Success)
	// 	p.SetString("return_msg", "ok")
	// 	p.SetString("data", xmlStr)
	// 	return p, err
	// }
}

//Report 交易保障
func (c *Client) Report(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = reportSandboxURL
	} else {
		url = reportURL
	}
	return sendAPI(key, url, params, false)
}

//ShortURL 转换短链接
func (c *Client) ShortURL(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = shortSandboxURL
	} else {
		url = shortURL
	}
	return sendAPI(key, url, params, false)
}

//AuthCodeToOpenID 授权码查询OPENID接口
func (c *Client) AuthCodeToOpenID(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = authCodeToOpenidSandboxURL
	} else {
		url = authCodeToOpenidURL
	}
	return sendAPI(key, url, params, false)
}
