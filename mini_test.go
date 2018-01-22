package main_test

import (
	"fmt"
	"github.com/beewit/wechat/mini"
	"testing"
)

var (
	appid  string
	secret string
	code   string

	encryptedData string
	iv            string
	sessionKey    string
)

func TestGetWxSessionKey(t *testing.T) {
	appid = "wx299f167a78297c8e"
	secret = "fdbacc37ad25d1b0db081a77dbdfb175"
	code = "003tDPrN0w4af52OYgrN0pBWrN0tDPr-"
	miniWx := mini.NewWx(appid, secret)
	ws, err := miniWx.GetWxSessionKey(code)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ws)
	if ws.ErrCode != 0 {
		t.Log(ws.ErrMsg)
		t.FailNow()
	}
	t.Log(ws.Openid)
	t.Log(ws.SessionKey)
	t.Log("0031i5Dl0mrmgl1PJECl0yc1Dl01i5DI")
	signature := `a83d982205e1be5a79fe39177dffe332d2ae1944`
	rawData := `{"nickName":"承诺，一时的华丽","gender":1,"language":"zh_CN","city":"Shapingba","province":"Chongqing","country":"China","avatarUrl":"https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKw8ictgYcqf6uklrSAup13EoCQ2SyfASGwOkOicAFibBz7LVgyPm7DoMPDhcqzNZzgXsJWt3r1l9gxQ/0"}`
	if mini.CheckSignature(signature, ws.SessionKey, rawData) != true {
		t.Fatal("CheckSignature failed")
		t.FailNow()
	}
	u, err := mini.GetWxUserInfo(sessionKey, encryptedData, iv)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(u)

}

func TestUser(t *testing.T) {

}
