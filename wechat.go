package main

import (
	"github.com/beewit/wechat/mp/message"
)

func main() {
	var tm message.TemplateMessage
	tm.ToUser = "oRDhB03JEmQcIhY_W9Q0dTq1e7-s"
	tm.URL = "http://www.9ee3.com/"
	tm.TemplateId="gTCuxY42MKftPVnbhbd8VYeY1-_TTPLRxenbdW17_xg"
	tm.Data = map[string]interface{}{
		"first":    "《图形验证码》输入提醒通知",
		"keyword1": map[string]string{"value": "1234", "color": "#173177"},
		"keyword2": map[string]string{"value": "10分钟过期", "color": "#173177"},
		"remark":   "value",
	}
	msgId, err := message.Send(tm,"")
	if err != nil {
		println("error:", err.Error())
	} else {
		println(msgId)
	}
}
