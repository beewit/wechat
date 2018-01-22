package mmpaymkttransfers

import (
	"github.com/beewit/wechat/mp"
)

// 发放裂变红包.
//  NOTE: 请求需要双向证书
func SendGroupRedPack(pxy *mp.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendgroupredpack", req)
}
