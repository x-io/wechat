package wechat

import (
	"github.com/x-io/wechat/config"
	"github.com/x-io/wechat/param"
)

//Params Params
type Params = param.Params

//Config Config
type Config = config.Config

//ConfigInit ConfigInit
type ConfigInit = func(key string) (*Config, error)
