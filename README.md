## 1. Install

```
go get -u github.com/x-io/wechat
```

## 2. Getting Started
```go
func main() {
    //cache.Init(&adapter.FreeCache{freecache.NewCache(1024)})
	wechat.Init(func(key string) (*wechat.Config, error) {
        if key == "Test" {
            config := new(wechat.Config)
            config.AppID = ""
            config.AppSecret = ""
            config.Pay.MchID = ""
            config.Pay.APIKey = ""
            //config.Pay.Cert, err = base64.StdEncoding.DecodeString(model.Pay.Cert)
            // if err != nil {
            // 	return nil, errors.New("证书数据为空")
            // }

            return config, nil
        }

		return nil, errors.New("未找到相应支付账号信息")
	})

	// 开启沙箱模式
    if err := wechat.Pay.Sandbox(); err != nil {
        return err
    }
    fmt.Println("微信支付 沙箱模式已启动")
	

    //创建支付订单
	params := make(wechat.Params)
    params.
        Set("body", "Demo").
        Set("openid", "").
        Set("total_fee", 100).
        Set("trade_type", "JSAPI").
        Set("spbill_create_ip", "127.0.0.1").
        Set("out_trade_no", "123456789").
        Set("notify_url", "http://xxxx")

    params, err = wechat.Pay.UnifiedOrder("Test", params)

    if err != nil {
        return "", err
    }

    prepayID = params.Get("prepay_id")

}
```
