package pay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/x-io/wechat/cache"
	"github.com/x-io/wechat/util"
)

//Sign 签名
func sign(signType, apiKey string, params Params) string {

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

	switch signType {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(apiKey))
		h.Write(buf.Bytes())
		dataSha256 = h.Sum(nil)
		str = hex.EncodeToString(dataSha256[:])
	default:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	}

	return strings.ToUpper(str)
}

//ValidSign 验证签名
func validSign(signType, apiKey string, params Params) bool {
	if apiKey == "" {
		return false
	}

	if !params.ContainsKey(Sign) {
		return false
	}

	return params.Get(Sign) == sign(signType, apiKey, params)
}

func getSandboxKey(config *cache.Config) (string, error) {
	url := getSignkeyURL
	h := &http.Client{}
	params := make(Params)
	params["nonce_str"] = util.NonceStr()
	params["appid"] = config.AppID
	params["mch_id"] = config.Pay.MchID
	params["sign_type"] = config.Pay.SignType
	params["sign"] = sign(config.Pay.SignType, config.Pay.APIKey, params)

	response, err := h.Post(url, "application/xml; charset=utf-8", strings.NewReader(util.MapToXML(params)))
	if err != nil {
		return "", err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	params, err = processXML(string(res))
	if err != nil {
		return "", err
	}

	if params.ContainsKey("sandbox_signkey") {
		return params.Get("sandbox_signkey"), nil
	}

	return "", errors.New("wechat: sandbox error")
}

func sendAPI(key string, url string, params Params, cert bool) (Params, error) {
	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	var h *http.Client
	if cert {
		if config.Pay.Transport == nil {
			var certData []byte

			if len(config.Pay.CertFile) > 0 {
				certData, err = ioutil.ReadFile(config.Pay.CertFile)
			}

			if len(certData) == 0 && len(config.Pay.CertData) > 0 {
				certData, err = base64.StdEncoding.DecodeString(config.Pay.CertData)
			}
			if len(certData) == 0 {
				return nil, errors.New("wechat: 证书不存在")
			}

			cert, err := util.P12ToPem(certData, config.Pay.MchID)
			if err != nil {
				return nil, errors.New("wechat: 证书解析无效")
			}

			config.Pay.Transport = &http.Transport{
				TLSClientConfig:    &tls.Config{Certificates: []tls.Certificate{cert}},
				DisableCompression: true,
			}
		}
		h = &http.Client{Transport: config.Pay.Transport}
	} else {
		h = &http.Client{}
	}

	apiKey := config.Pay.APIKey
	if config.Pay.SandboxKey != "" {
		apiKey = config.Pay.SandboxKey
	} else {
		if strings.Index(url, "/sandboxnew/") > 0 {
			signkey, err := getSandboxKey(config)
			if err != nil {
				return nil, err
			}
			apiKey = signkey
			config.Pay.SandboxKey = signkey
		}
	}

	params["appid"] = config.AppID
	params["mch_id"] = config.Pay.MchID
	params["sign_type"] = config.Pay.SignType
	params["nonce_str"] = util.NonceStr()
	params["sign"] = sign(config.Pay.SignType, apiKey, params)

	response, err := h.Post(url, "application/xml; charset=utf-8", strings.NewReader(util.MapToXML(params)))
	if err != nil {
		return nil, err
	}

	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	params, err = processXML(string(res))
	if err != nil {
		return nil, err
	}

	if validSign(config.Pay.SignType, apiKey, params) {
		return params, nil
	}

	return nil, errors.New("wechat: invalid sign error")
}

func sendTransfer(key string, url string, params Params, cert bool) (Params, error) {
	config, err := cache.GetConfig(key)
	if err != nil {
		return nil, err
	}

	var h *http.Client
	if cert {
		if config.Pay.Transport == nil {
			var certData []byte

			if len(config.Pay.CertFile) > 0 {
				certData, err = ioutil.ReadFile(config.Pay.CertFile)
			}

			if len(certData) == 0 && len(config.Pay.CertData) > 0 {
				certData, err = base64.StdEncoding.DecodeString(config.Pay.CertData)
			}
			if len(certData) == 0 {
				return nil, errors.New("wechat: 证书不存在")
			}

			cert, err := util.P12ToPem(certData, config.Pay.MchID)
			if err != nil {
				return nil, errors.New("wechat: 证书解析无效")
			}

			config.Pay.Transport = &http.Transport{
				TLSClientConfig:    &tls.Config{Certificates: []tls.Certificate{cert}},
				DisableCompression: true,
			}
		}
		h = &http.Client{Transport: config.Pay.Transport}
	} else {
		h = &http.Client{}
	}

	params["mch_appid"] = config.AppID
	params["mchid"] = config.Pay.MchID
	params["sign_type"] = config.Pay.SignType
	params["nonce_str"] = util.NonceStr()
	params["sign"] = sign(config.Pay.SignType, config.Pay.APIKey, params)

	response, err := h.Post(url, "application/xml; charset=utf-8", strings.NewReader(util.MapToXML(params)))
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return processXML(string(res))
}

// 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func processXML(xmlStr string) (Params, error) {
	var returnCode string
	params := Params(util.XMLToMap(xmlStr))

	if !params.ContainsKey("return_code") {
		return nil, errors.New("wechat: 远程微信服务器异常")
	}

	returnCode = params.Get("return_code")

	if returnCode == Success {
		if params.Get("result_code") == Fail {
			return params, errors.New("wechat: " + params.Get("err_code_des"))
		}

		return params, nil
	} else if returnCode == Fail {
		if params.Get("return_msg") != "" {
			return params, errors.New("wechat: " + params.Get("return_msg"))
		} else if params.Get("retmsg") != "" {
			return params, errors.New("wechat: " + params.Get("retmsg"))
		}
		return params, errors.New("wechat: 参数错误")
	}

	return nil, errors.New("wechat: return_code value is invalid in XML")
}
