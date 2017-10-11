package global

import (
	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/redis"
	"github.com/beewit/beekit/log"
	"github.com/astaxie/beego/logs"
)

var (
	CFG conf.Config
	RD  *redis.RedisConnPool
	App *AppWechat
	Log *logs.BeeLogger
)

func init() {
	CFG = conf.New("config.json")
	RD = redis.Cache
	Log = log.Logger
	App = &AppWechat{
		AppID:     convert.ToString(CFG.Get("wechat.appId")),
		AppSecret: convert.ToString(CFG.Get("wechat.appSecret")),
		MchID:     convert.ToString(CFG.Get("wechat.mchID")),
		APIKey:    convert.ToString(CFG.Get("wechat.apiKey")),
	}
}

type AppWechat struct {
	AppID          string
	MchID          string
	AppSecret      string
	APIKey         string
	EncodingAESKey string
}
