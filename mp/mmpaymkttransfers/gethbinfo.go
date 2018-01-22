package mmpaymkttransfers

import (
	"github.com/beewit/wechat/mp"
)

// 红包查询接口.
//  NOTE: 请求需要双向证书
func GetRedPackInfo(pxy *mp.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/mmpaymkttransfers/gethbinfo", req)
}
