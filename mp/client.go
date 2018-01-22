package mp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beewit/wechat/util"
	"net/http"
	"net/url"
	"reflect"
)

type Client struct {
	AccessTokenServer
	HttpClient *http.Client
}

// 创建一个新的 Client.
//  如果 clt == nil 则默认用 http.DefaultClient
func NewClient(srv AccessTokenServer, clt *http.Client) *Client {
	if srv == nil {
		panic("nil AccessTokenServer")
	}
	if clt == nil {
		clt = http.DefaultClient
	}

	return &Client{
		AccessTokenServer: srv,
		HttpClient:        clt,
	}
}

// 用 encoding/json 把 request marshal 为 JSON, 放入 http 请求的 body 中,
// POST 到微信服务器, 然后将微信服务器返回的 JSON 用 encoding/json 解析到 response.
//
//  NOTE:
//  1. 一般不用调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + access_token;
//  3. response 格式有要求, 要么是 *Error, 要么是下面结构体的指针(注意 Error 必须是第一个 Field):
//      struct {
//          Error
//          ...
//      }
func (clt *Client) PostJSON(incompleteURL string, request interface{}, response interface{}) (err error) {
	buf := textBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer textBufferPool.Put(buf)

	if err = json.NewEncoder(buf).Encode(request); err != nil {
		return
	}
	requestBytes := buf.Bytes()

	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Post(finalURL, "application/json; charset=utf-8", bytes.NewReader(requestBytes))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	if err = json.NewDecoder(httpResp.Body).Decode(response); err != nil {
		return
	}

	var ErrorStructValue reflect.Value // Error

	// 下面的代码对 response 有特定要求, 见此函数 NOTE
	responseStructValue := reflect.ValueOf(response).Elem()
	if v := responseStructValue.Field(0); v.Kind() == reflect.Struct {
		ErrorStructValue = v
	} else {
		ErrorStructValue = responseStructValue
	}

	switch ErrCode := ErrorStructValue.Field(0).Int(); ErrCode {
	case ErrCodeOK:
		return
	case ErrCodeInvalidCredential, ErrCodeAccessTokenExpired:
		ErrMsg := ErrorStructValue.Field(1).String()
		util.LogInfoln("[WECHAT_RETRY] err_code:", ErrCode, ", err_msg:", ErrMsg)
		util.LogInfoln("[WECHAT_RETRY] current token:", token)

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			util.LogInfoln("[WECHAT_RETRY] new token:", token)

			responseStructValue.Set(reflect.New(responseStructValue.Type()).Elem())
			goto RETRY
		}
		util.LogInfoln("[WECHAT_RETRY] fallthrough, current token:", token)
		fallthrough
	default:
		return
	}
}

// GET 微信资源, 然后将微信服务器返回的 JSON 用 encoding/json 解析到 response.
//
//  NOTE:
//  1. 一般不用调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + access_token;
//  3. response 格式有要求, 要么是 *Error, 要么是下面结构体的指针(注意 Error 必须是第一个 Field):
//      struct {
//          Error
//          ...
//      }
func (clt *Client) GetJSON(incompleteURL string, response interface{}) (err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	if err = json.NewDecoder(httpResp.Body).Decode(response); err != nil {
		return
	}

	var ErrorStructValue reflect.Value // Error

	// 下面的代码对 response 有特定要求, 见此函数 NOTE
	responseStructValue := reflect.ValueOf(response).Elem()
	if v := responseStructValue.Field(0); v.Kind() == reflect.Struct {
		ErrorStructValue = v
	} else {
		ErrorStructValue = responseStructValue
	}

	switch ErrCode := ErrorStructValue.Field(0).Int(); ErrCode {
	case ErrCodeOK:
		return
	case ErrCodeInvalidCredential, ErrCodeAccessTokenExpired:
		ErrMsg := ErrorStructValue.Field(1).String()
		util.LogInfoln("[WECHAT_RETRY] err_code:", ErrCode, ", err_msg:", ErrMsg)
		util.LogInfoln("[WECHAT_RETRY] current token:", token)

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			util.LogInfoln("[WECHAT_RETRY] new token:", token)

			responseStructValue.Set(reflect.New(responseStructValue.Type()).Elem())
			goto RETRY
		}
		util.LogInfoln("[WECHAT_RETRY] fallthrough, current token:", token)
		fallthrough
	default:
		return
	}
}

type Proxy struct {
	appId      string
	mchId      string
	apiKey     string
	httpClient *http.Client
}

func (pxy *Proxy) AppId() string {
	return pxy.appId
}
func (pxy *Proxy) MchId() string {
	return pxy.mchId
}

// 创建一个新的 Proxy.
//  如果 httpClient == nil 则默认用 http.DefaultClient.
func NewProxy(appId, mchId, apiKey string, httpClient *http.Client) *Proxy {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Proxy{
		appId:      appId,
		mchId:      mchId,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// 微信支付通用请求方法.
//  注意: err == nil 表示协议状态都为 SUCCESS(return_code == SUCCESS).
func (pxy *Proxy) PostXML(url string, req map[string]string) (resp map[string]string, err error) {
	bodyBuf := textBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer textBufferPool.Put(bodyBuf)

	if err = util.EncodeXMLFromMap(bodyBuf, req, "xml"); err != nil {
		return
	}

	httpResp, err := pxy.httpClient.Post(url, "text/xml; charset=utf-8", bodyBuf)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	if resp, err = util.DecodeXMLToMap(httpResp.Body); err != nil {
		return
	}

	bs, _ := json.Marshal(resp)
	respJson := string(bs)

	// 判断协议状态
	ReturnCode, ok := resp["return_code"]
	if !ok {
		err = fmt.Errorf("%s ,no return_code parameter", respJson)
		return
	}
	if ReturnCode != ReturnCodeSuccess {
		err = &Error2{
			ReturnCode: ReturnCode,
			ReturnMsg:  resp["return_msg"],
		}
		return
	}

	//判断是否错误
	resultCode, ok := resp["result_code"]
	if !ok {
		err = fmt.Errorf("%s ,no return_code parameter", respJson)
		return
	}
	if resultCode != ResultCodeSuccess {
		err = fmt.Errorf("%s ,result_code FAIL, err_code：%s, err_code_des：%s", respJson, resp["err_code"], resp["err_code_des"])
		return
	}

	// 安全考虑, 做下验证
	appId, ok := resp["appid"]
	if ok && appId != pxy.appId {
		err = fmt.Errorf("%s ,appid mismatch, have: %q, want: %q", respJson, appId, pxy.appId)
		return
	}
	mchId, ok := resp["mch_id"]
	if ok && mchId != pxy.mchId {
		err = fmt.Errorf("%s ,mch_id mismatch, have: %q, want: %q", respJson, mchId, pxy.mchId)
		return
	}

	// 认证签名
	signature1, ok := resp["sign"]
	if !ok {
		err = fmt.Errorf("%s ,no sign parameter", respJson)
		return
	}
	signature2 := Sign(resp, pxy.apiKey, nil)
	if signature1 != signature2 {
		err = fmt.Errorf("%s ,check signature failed, \r\ninput: %q, \r\nlocal: %q", respJson, signature1, signature2)
		return
	}
	return
}
