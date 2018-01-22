package mp

import (
	"net/url"
	"net/http"
)

// http 请求拦截器
type Interceptor interface {
	// 拦截 http 请求, 根据需要做一些判断, 返回是否允许后续逻辑继续处理请求, 如返回 false 则表示请求到此为止.
	// 请注意, 后续逻辑需要读取 r.Body 里的内容, 请谨慎读取!
	Intercept(w http.ResponseWriter, r *http.Request, queryValues url.Values) (shouldContinue bool)
}

type InterceptorFunc func(w http.ResponseWriter, r *http.Request, queryValues url.Values) (shouldContinue bool)

func (fn InterceptorFunc) Intercept(w http.ResponseWriter, r *http.Request, queryValues url.Values) (shouldContinue bool) {
	return fn(w, r, queryValues)
}
