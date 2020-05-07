package wechat

import (
	"github.com/x-io/wechat/api"
	"github.com/x-io/wechat/config"
	"github.com/x-io/wechat/oauth"
	"github.com/x-io/wechat/pay"
)

var (
	API   *api.API
	OAuth *oauth.OAuth
	Pay   *pay.Client
)

//Init Init
func Init(call ConfigInit) {
	config.SetConfig(call)
	API = api.New()
	Pay = pay.New()
	OAuth = oauth.New()
}
