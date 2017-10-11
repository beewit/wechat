package util

const (
	ErrCodeOK                 = 0
	ErrCodeInvalidCredential  = 40001 // access_token 过期错误码
	ErrCodeAccessTokenExpired = 42001 // access_token 过期错误码(maybe!!!)
)

//微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
