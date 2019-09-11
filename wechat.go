package wechat

import (
	"github.com/x-io/wechat/api"
	"github.com/x-io/wechat/cache"
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
	cache.SetConfig(call)
	API = api.New()
	OAuth = oauth.New()
	Pay = pay.New()
}
