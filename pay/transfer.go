package pay

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/x-io/wechat/cache"
	"github.com/x-io/wechat/util"
)

// 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func (c *Client) processXML2(xmlStr string) (Params, error) {
	var returnCode string
	params := Params(util.XMLToMap(xmlStr))

	if params.ContainsKey("return_code") {
		returnCode = params.Get("return_code")
	} else {
		return nil, errors.New("接口返回异常")
	}

	if returnCode == Success {
		if params.Get("result_code") == Success {
			return params, nil
		}
		return params, errors.New(params.Get("err_code_des"))
	} else if returnCode == Fail {
		return params, errors.New(params.Get("return_msg"))
	} else {
		return nil, errors.New(params.Get("return_msg"))
	}
}

func (c *Client) sendTransfer(key string, url string, params Params, cert bool) (Params, error) {
	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	var h *http.Client
	if cert {
		if config.Pay.Transport == nil {
			transport, _ := util.GetTransport(config.Pay.Cert, config.Pay.MchID)

			if transport == nil {
				return nil, errors.New("证书数据为空")
			}
			config.Pay.Transport = transport
		}
		h = &http.Client{Transport: config.Pay.Transport}
	} else {
		h = &http.Client{}
	}

	params["mch_appid"] = config.AppID
	params["mchid"] = config.Pay.MchID
	params["nonce_str"] = util.NonceStr()
	//params["sign_type"] = c.signType
	params["sign"] = c.sign(config.Pay.APIKey, params)

	response, err := h.Post(url, "application/xml; charset=utf-8", strings.NewReader(util.MapToXML(params)))
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return c.processXML2(string(res))
}

//Transfer 转账到钱包
func (c *Client) Transfer(key string, params Params) (Params, error) {
	if c.isSandbox {
		return nil, errors.New("该方法不支持沙盒模式")
	}
	return c.sendTransfer(key, transferURL, params, true)
}

//TransferBank 转账到钱包
func (c *Client) TransferBank(key string, params Params) (Params, error) {
	if c.isSandbox {
		return nil, errors.New("该方法不支持沙盒模式")
	}
	return c.sendTransfer(key, transferBankURL, params, true)
}
