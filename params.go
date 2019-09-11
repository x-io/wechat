package wechat

import (
	"github.com/x-io/wechat/cache"
	"github.com/x-io/wechat/param"
)

//Params Params
type Params = param.Params

//Config Config
type Config = cache.Config

//ConfigInit ConfigInit
type ConfigInit = func(key string) (*Config, error)
