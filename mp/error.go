package mp

import "fmt"

const (
	ErrCodeOK                 = 0
	ErrCodeInvalidCredential  = 40001 // access_token 过期(无效)返回这个错误
	ErrCodeAccessTokenExpired = 42001 // access_token 过期(无效)返回这个错误(maybe!!!)
)

type Error struct {
	// NOTE: StructField 固定这个顺序, RETRY 依赖这个顺序
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
}

type Error2 struct {
	XMLName    struct{} `xml:"xml"                  json:"-"`
	ReturnCode string   `xml:"return_code"          json:"return_code"`
	ReturnMsg  string   `xml:"return_msg,omitempty" json:"return_msg,omitempty"`
}

func (e *Error2) Error() string {
	return fmt.Sprintf("return_code: %q, return_msg: %q", e.ReturnCode, e.ReturnMsg)
}
