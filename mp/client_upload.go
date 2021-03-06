package mp

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"

	"encoding/json"
)

type MultipartFormField struct {
	ContentType int // 0:文件field, 1:普通的文本field
	FieldName   string
	FileName    string // ContentType == 0 时有效
	Value       io.Reader
}

// 通用上传接口.
//
//  --BOUNDARY
//  Content-Disposition: form-data; name="FIELDNAME"; filename="FILENAME"
//  Content-Type: application/octet-stream
//
//  FILE-CONTENT
//  --BOUNDARY
//  Content-Disposition: form-data; name="FIELDNAME"
//
//  JSON-DESCRIPTION
//  --BOUNDARY--
//
//
//  NOTE:
//  1. 一般不需要调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + access_token;
//  3. response 格式有要求, 要么是 *Error, 要么是下面结构体的指针(注意 Error 必须是第一个 Field):
//      struct {
//          Error
//          ...
//      }
func (clt *Client) PostMultipartForm(incompleteURL string, fields []MultipartFormField, response interface{}) (err error) {
	bodyBuf := mediaBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer mediaBufferPool.Put(bodyBuf)

	multipartWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		switch field.ContentType {
		case 0: // 文件
			partWriter, err := multipartWriter.CreateFormFile(field.FieldName, field.FileName)
			if err != nil {
				return err
			}
			if _, err = io.Copy(partWriter, field.Value); err != nil {
				return err
			}
		case 1: // 文本
			partWriter, err := multipartWriter.CreateFormField(field.FieldName)
			if err != nil {
				return err
			}
			if _, err = io.Copy(partWriter, field.Value); err != nil {
				return err
			}
		}
	}

	if err = multipartWriter.Close(); err != nil {
		return
	}

	bodyBytes := bodyBuf.Bytes()

	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Post(finalURL, multipartWriter.FormDataContentType(), bytes.NewReader(bodyBytes))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
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
		//ErrMsg := ErrorStructValue.Field(1).String()

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}

			responseStructValue.Set(reflect.New(responseStructValue.Type()).Elem())
			goto RETRY
		}
		fallthrough
	default:
		return
	}
}
