package oauth2

import (
	"github.com/beewit/wechat/util"
	"github.com/beewit/wechat/global"
	"encoding/json"
	"errors"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

func GetAccessToken(code string) *AccessToken {
	var result struct {
		util.CommonError
		AccessToken
	}
	b, err := util.PostJSON(util.GetAuthorizeAccessTokenUrl(code), nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if result.ErrCode != util.ErrCodeOK {
		err = errors.New(result.ErrMsg)
		return nil
	}
	return &result.AccessToken
}
