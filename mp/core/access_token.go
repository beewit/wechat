package core

import (
	"fmt"
	"github.com/beewit/wechat/util"
	"github.com/beewit/wechat/global"
	"encoding/json"
	"sync"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
)

var AccessToken *AppAccessToken

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
	util.CommonError
}

type AppAccessToken struct {
	*global.AppWechat
	accessTokenLock *sync.RWMutex //accessTokenLock 读写锁 同一个AppID一个
	jsAPITicketLock *sync.RWMutex //jsAPITicket 读写锁 同一个AppID一个
}

func init() {
	AccessToken = &AppAccessToken{}
	AccessToken.AppWechat = global.App
	AccessToken.accessTokenLock = new(sync.RWMutex)
	AccessToken.jsAPITicketLock = new(sync.RWMutex)
}

//GetAccessToken 获取access_token
func (app *AppAccessToken) GetAccessToken() (accessToken string, err error) {
	app.accessTokenLock.Lock()
	defer app.accessTokenLock.Unlock()

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", app.AppID)
	val, _ := global.RD.GetString(accessTokenCacheKey)
	if val != "" {
		accessToken = val
		return
	}
	resAccessToken, err := app.GetAccessTokenFromServer()
	if err != nil {
		return
	}
	accessToken = resAccessToken.Token
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (app *AppAccessToken) GetAccessTokenFromServer() (resAccessToken accessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, app.AppID, app.AppSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	global.Log.Info(string(body))
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", app.AppID)
	expires := resAccessToken.ExpiresIn - 1500
	_, err = global.RD.SetAndExpire(accessTokenCacheKey, resAccessToken.Token, expires)

	return
}
