package pay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/x-io/wechat/cache"
	"github.com/x-io/wechat/param"
	"github.com/x-io/wechat/util"
)

//Params Params
type Params = param.Params

//Client Pay Client
type Client struct {
	isSandbox bool
	signType  string
}

//New 创建微信支付客户端
func New() *Client {
	return &Client{signType: MD5}
}

// 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func (c *Client) processXML(apiKey string, xmlStr string) (Params, error) {
	var returnCode string
	params := Params(util.XMLToMap(xmlStr))
	// fmt.Println(xmlStr, params)
	if params.ContainsKey("return_code") {
		returnCode = params.Get("return_code")
	} else {
		return nil, errors.New("远程微信服务器异常")
	}

	if returnCode == Fail {
		return params, errors.New(params.Get("return_msg"))
	} else if returnCode == Success {
		if c.validSign(apiKey, params) {
			return params, nil
		}
		if c.isSandbox && params.ContainsKey("sandbox_signkey") {
			return params, nil
		}
		//fmt.Println(apiKey,params)
		return nil, errors.New("invalid sign value in XML")
	} else {
		return nil, errors.New("return_code value is invalid in XML")
	}
}

func (c *Client) send(key string, url string, params Params, cert bool) (Params, error) {
	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	var h *http.Client
	if cert {
		if config.Pay.Transport == nil {
			config.Pay.Transport, err = util.GetTransport(config.Pay.Cert, config.Pay.MchID)
			if err != nil {
				return nil, errors.New("证书数据为空")
			}
		}
		h = &http.Client{Transport: config.Pay.Transport}
	} else {
		h = &http.Client{}
	}

	params["appid"] = config.AppID
	params["mch_id"] = config.Pay.MchID
	params["sign_type"] = c.signType
	params["nonce_str"] = util.NonceStr()
	if c.isSandbox {
		if config.Pay.SandboxKey == "" {
			signkey, err := c.getSandboxKey(config)
			if err != nil {
				return nil, err
			}
			config.Pay.SandboxKey = signkey
		}
		params["sign"] = c.sign(config.Pay.SandboxKey, params)

	} else {
		params["sign"] = c.sign(config.Pay.APIKey, params)
	}

	response, err := h.Post(url, "application/xml; charset=utf-8", strings.NewReader(util.MapToXML(params)))
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if c.isSandbox {
		return c.processXML(config.Pay.SandboxKey, string(res))
	}
	return c.processXML(config.Pay.APIKey, string(res))
}

//Sign 签名
func (c *Client) sign(apiKey string, params Params) string {

	// 创建切片
	var keys = make([]string, 0, len(params))
	// 遍历签名参数
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}

	// 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
	sort.Strings(keys)

	//创建字符缓冲
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.Get(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.Get(k))
			buf.WriteString(`&`)
		}
	}
	// 加入apiKey作加密密钥
	buf.WriteString(`key=`)
	buf.WriteString(apiKey)

	var (
		dataMd5    [16]byte
		dataSha256 []byte
		str        string
	)

	switch c.signType {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(apiKey))
		h.Write(buf.Bytes())
		dataSha256 = h.Sum(nil)
		str = hex.EncodeToString(dataSha256[:])
	}

	return strings.ToUpper(str)
}

//ValidSign 验证签名
func (c *Client) validSign(apiKey string, params Params) bool {
	//fmt.Println(apiKey, params, c.sign(apiKey, params))
	if apiKey == "" {
		return false
	}

	if !params.ContainsKey(Sign) {
		return false
	}

	return params.Get(Sign) == c.sign(apiKey, params)
}

//Sandbox 开启沙箱模式
func (c *Client) Sandbox() error {
	c.isSandbox = true
	return nil
}

func (c *Client) getSandboxKey(config *cache.Config) (string, error) {
	url := getSignkeyURL
	h := &http.Client{}
	params := make(Params)
	params["sign_type"] = c.signType
	params["nonce_str"] = util.NonceStr()
	params["appid"] = config.AppID
	params["mch_id"] = config.Pay.MchID
	params["sign"] = c.sign(config.Pay.APIKey, params)

	response, err := h.Post(url, "application/xml; charset=utf-8", strings.NewReader(util.MapToXML(params)))
	if err != nil {
		return "", err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	params, err = c.processXML(config.Pay.APIKey, string(res))
	if err != nil {
		return "", err
	}

	return params.Get("sandbox_signkey"), nil
}

//ValidNotify 交易通知验证
func (c *Client) ValidNotify(body string) (Params, error) {

	params := Params(util.XMLToMap(body))

	config, err := cache.GetConfig(params.Get("attach"))
	if err != nil {
		return nil, err
	}

	if c.validSign(config.Pay.APIKey, params) {
		return params, nil
	}

	return params, errors.New("签名认证失败")
}

// //ValidNotify 交易通知验证
// func (c *Client) ValidNotify(key string, body string) (Params, error) {

// 	config, err := cache.GetConfig(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	params := Params(util.XMLToMap(body))

// 	if c.validSign(config.Pay.APIKey, params) {
// 		return params, nil
// 	}

// 	return params, errors.New("签名认证失败")
// }

//ChooseWXPay 获取Js调用参数
func (c *Client) ChooseWXPay(key string, prepayID string) (Params, error) {
	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	params := make(Params).Set("timeStamp", util.GetCurrTs())
	params["appId"] = config.AppID
	params["nonceStr"] = util.NonceStr()
	params["signType"] = c.signType
	params["package"] = "prepay_id=" + prepayID
	params["paySign"] = c.sign(config.Pay.APIKey, params)
	params["timestamp"] = params.Get("timeStamp")
	delete(params, "timeStamp")
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
	return c.send(key, url, params.Set("attach", key), false)
}

//MicroPay 刷卡支付
func (c *Client) MicroPay(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = microPaySandboxURL
	} else {
		url = microPayURL
	}
	return c.send(key, url, params, false)
}

//Refund 退款
func (c *Client) Refund(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = refundSandboxURL
	} else {
		url = refundURL
	}
	return c.send(key, url, params, true)
}

//OrderQuery 订单查询
func (c *Client) OrderQuery(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = orderQuerySandboxURL
	} else {
		url = orderQueryURL
	}
	return c.send(key, url, params, false)
}

//RefundQuery 退款查询
func (c *Client) RefundQuery(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = refundQuerySandboxURL
	} else {
		url = refundQueryURL
	}
	return c.send(key, url, params, false)
}

//Reverse 撤销订单
func (c *Client) Reverse(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = reverseSandboxURL
	} else {
		url = reverseURL
	}
	return c.send(key, url, params, true)
}

//CloseOrder 关闭订单
func (c *Client) CloseOrder(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = closeOrderSandboxURL
	} else {
		url = closeOrderURL
	}
	return c.send(key, url, params, false)
}

//DownloadBill 对账单下载
func (c *Client) DownloadBill(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = downloadBillSandboxURL
	} else {
		url = downloadBillURL
	}
	return c.send(key, url, params, false)

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
	return c.send(key, url, params, false)
}

//ShortURL 转换短链接
func (c *Client) ShortURL(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = shortSandboxURL
	} else {
		url = shortURL
	}
	return c.send(key, url, params, false)
}

//AuthCodeToOpenid 授权码查询OPENID接口
func (c *Client) AuthCodeToOpenid(key string, params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = authCodeToOpenidSandboxURL
	} else {
		url = authCodeToOpenidURL
	}
	return c.send(key, url, params, false)
}
