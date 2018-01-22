package mini

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beewit/wechat/util"
)

type Wx struct {
	AppId     string
	AppSecret string
}

//微信加密数据结构
type WxUserInfo struct {
	OpenId    string     `json:"openId"`
	NickName  string     `json:"nickName"`
	Gender    int        `json:"gender"`
	City      string     `json:"city"`
	Province  string     `json:"province"`
	Country   string     `json:"country"`
	AvatarUrl string     `json:"avatarUrl"`
	UnionId   string     `json:"unionId"`
	Watermark *Watermark `json:"watermark"` //数据水印( watermark )
}
type Watermark struct {
	Appid     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}
type WxSesstion struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrInfo
}
type ErrInfo struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewWx(appid, secret string) *Wx {
	return &Wx{
		AppId:     appid,
		AppSecret: secret,
	}
}

//"https://api.weixin.qq.com/sns/jscode2session?
//appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
//根据code获取WxSesstion
func (wx *Wx) GetWxSessionKey(code string) (ws WxSesstion, err error) {
	uri := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?"+
		"appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		wx.AppId, wx.AppSecret, code)
	err = util.HttpGet(uri).ToJson(&ws)
	return
}

//检验signature是否相同
func CheckSignature(signature, session_key, rawData string) bool {
	if signature == Sha1(rawData+session_key) {
		return true
	}
	return false
}

//根据seesion_key,加密数据encryptedData和向量偏移量iv获取微信用户信息 主要是 敏感信息
func GetWxUserInfo(session_key, encryptedData, iv string) (wui WxUserInfo, err error) {
	plaintext, err := AesCBCDecrypt(session_key, encryptedData, iv)
	if err != nil {
		return
	}
	err = json.Unmarshal(plaintext, &wui)
	return
}

func Sha1(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

var (
	ErrPaddingSize = errors.New("padding size error")
)

// 接口返回的加密数据( encryptedData )进行对称解密。 解密算法如下：
// 对称解密使用的算法为 AES-128-CBC，数据采用PKCS#7填充。
// 对称解密的目标密文为 Base64_Decode(encryptedData)。
// 对称解密秘钥 aeskey = Base64_Decode(session_key), aeskey 是16字节。
// 对称解密算法初始向量 为Base64_Decode(iv)，其中iv由数据接口返回。
// AES并没有64位的块, 如果采用PKCS5, 那么实质上就是采用PKCS7
func AesCBCDecrypt(session_key, encryptedData, iv string) (plaintext []byte,
	err error) {
	// base
	aeskey, err := base64.StdEncoding.DecodeString(session_key)
	if err != nil {
		return nil, errors.New("base64 decoding session_key err:" + err.Error())
	}
	encryptedDatabytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.New("base64 decoding encryptedData err:" + err.Error())
	}
	ivbytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.New("base64 decoding iv err:" + err.Error())
	}
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	plaintext = make([]byte, len(encryptedDatabytes))
	cipher.NewCBCDecrypter(block, ivbytes).CryptBlocks(plaintext, encryptedDatabytes)
	return PKCS5UnPadding(plaintext, block.BlockSize())
}

func PKCS5UnPadding(plaintext []byte, blockSize int) ([]byte, error) {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	//
	if unpadding >= length || unpadding > blockSize {
		return nil, ErrPaddingSize
	}
	return plaintext[:(length - unpadding)], nil
}
