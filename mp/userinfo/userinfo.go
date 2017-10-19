package userinfo

import (
	"encoding/json"
	"errors"
	"github.com/beewit/wechat/global"
	"github.com/beewit/wechat/util"
)

/**
{"subscribe":1,"openid":"oRDhB09D3al5ICDjCMfy3ac1DRP0","nickname":"执手并肩看天下","sex":1,"language":"zh_CN","city":"","province":"重庆","country":"中国","headimgurl":"http:\/\/wx.qlogo.cn\/mmopen\/BicliabeoVxSEEAibcy15aw6LNtXI0NENAsTLy40nd2ZMUOxgOHkkAicx09DibSR5HCP1H4FQTYx2msbCy2ibtWmdrj2UeianacMfel\/0","subscribe_time":1507798056,"unionid":"oWYCdv-DoLpnOHjx3gnPSIA3tvaU","remark":"","groupid":0,"tagid_list":[]}
/
*/
type UserInfo struct {
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
	Privilege  string `json:"privilege"`
	Unionid    string `json:"unionid"`
}

func GetUserInfo(accessToken, openId string) (*UserInfo, error) {
	var result struct {
		util.CommonError
		UserInfo
	}
	global.Log.Info("GetUserInfoUrl：%s", util.GetUserInfoUrl(accessToken, openId))
	b, err := util.PostJSON(util.GetUserInfoUrl(accessToken, openId), nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	global.Log.Info("GetUserInfo：%s", string(b))
	err = json.Unmarshal(b, &result)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	if result.ErrCode != util.ErrCodeOK {
		err = errors.New(result.ErrMsg)
		return nil, err
	}
	return &result.UserInfo, nil
}
