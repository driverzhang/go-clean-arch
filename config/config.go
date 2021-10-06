package config

import (
	ginConfig "git.dustess.com/mk-base/gin-ext/config"
	ginConstant "git.dustess.com/mk-base/gin-ext/constant"
	"git.dustess.com/mk-base/gonfig"
)

// Config 配置结构
type Config struct {
	ConfigFile string
	ginConfig.Config
	MKMongoConfig ginConfig.MKMongoConfig
	WPBillAddr    string `id:"wp-bill-addr" default:"" desc:"计费服务新兼容接口"`
	MaxPoolSize   struct {
		Markting int `id:"markting" default:"100" desc:"最大连接数设置, 每开pod就叠加"`
	} `id:"max_pool_size" desc:"最大连接数设置"`
	DigitalScreen string `id:"digital_screen" default:"http://82.157.33.81:31601/blade-visual/auth/authorization" desc:"数字大屏通知地址"`
}

// NewConfig 初始化配置
func NewConfig() *Config {
	config := &Config{}
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := config.MKMongoConfig.Init(); err != nil {
		panic(err)
	}
	err := gonfig.Load(config, gonfig.Conf{
		ConfigFileVariable:  ginConstant.ConfigFileVariable, // enables passing --configfile myfile.conf
		FileDefaultFilename: ginConstant.ConfigFileDefaultName,
		FileDecoder:         gonfig.DecoderJSON,
		EnvPrefix:           ginConstant.ConfigEnvPrefix,
	})
	if err != nil {
		panic(err)
	}
	return config
}
