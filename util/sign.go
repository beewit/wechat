package util

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"hash"
	"crypto/md5"
	"bytes"
)


// 微信支付签名.
//  parameters: 待签名的参数集合
//  apiKey:     API密钥
//  fn:         func() hash.Hash, 如果 fn == nil 则默认用 md5.New
func Sign(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	if fn == nil {
		fn = md5.New
	}
	h := fn()

	buf := make([]byte, 256)
	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}

		buf = buf[:0]
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, v...)
		buf = append(buf, '&')
		h.Write(buf)
	}
	buf = buf[:0]
	buf = append(buf, "key="...)
	buf = append(buf, apiKey...)
	h.Write(buf)

	signature := make([]byte, h.Size()*2)
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

// 微信公众号 明文模式/URL认证 签名
func UrlSign(token, timestamp, nonce string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce))

	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

// 微信公众号/企业号 密文模式消息签名
func MsgSign(token, timestamp, nonce, encryptedMsg string) (signature string) {
	strs := sort.StringSlice{token, timestamp, nonce, encryptedMsg}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce)+len(encryptedMsg))

	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)
	buf = append(buf, strs[3]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}
