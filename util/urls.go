package util

import (
	"fmt"
	"github.com/beewit/wechat/global"
	"net/url"
)

const (
	AUTHORIZE_CODE         = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect"
	AUTHORIZE_ACCESS_TOKEN = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	USERINFO               = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
)

func GetAuthorizeCodeUrl(redirectUri string) string {
	return fmt.Sprintf(AUTHORIZE_CODE, global.App.AppID, url.QueryEscape(redirectUri))
}

func GetAuthorizeAccessTokenUrl(code string) string {
	return fmt.Sprintf(AUTHORIZE_ACCESS_TOKEN, global.App.AppID, global.App.AppSecret, code)
}

func GetUserInfoUrl(accessToken, openId string) string {
	return fmt.Sprintf(USERINFO, accessToken, openId)
}
