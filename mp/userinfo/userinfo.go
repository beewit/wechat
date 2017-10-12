package userinfo

import (
	"github.com/beewit/wechat/util"
	"github.com/beewit/wechat/global"
	"encoding/json"
	"errors"
)

type UserInfo struct {
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        string `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
	Privilege  string `json:"privilege"`
	Unionid    string `json:"unionid"`
}

func GetUserInfo(accessToken, openId string) *UserInfo {
	var result struct {
		util.CommonError
		UserInfo
	}
	b, err := util.PostJSON(util.GetUserInfoUrl(accessToken, openId), nil)
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
	return &result.UserInfo
}
