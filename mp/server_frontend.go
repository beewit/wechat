package mp

import (
	"github.com/labstack/echo"
	"net/url"
)

// ServerFrontend 实现了 http.Handler, 处理一个公众号的消息(事件)请求.
type ServerFrontend struct {
	server      Server
	errHandler  ErrorHandler
	interceptor Interceptor
}

// NOTE: errHandler, interceptor 均可以为 nil
func NewServerFrontend(server Server, errHandler ErrorHandler, interceptor Interceptor) *ServerFrontend {
	if server == nil {
		panic("nil Server")
	}
	if errHandler == nil {
		errHandler = DefaultErrorHandler
	}

	return &ServerFrontend{
		server:      server,
		errHandler:  errHandler,
		interceptor: interceptor,
	}
}

func (frontend *ServerFrontend) ServeHTTP(c echo.Context) error {
	queryValues, err := url.ParseQuery(c.Request().URL.RawQuery)
	if err != nil {
		frontend.errHandler.ServeError(c.Response(), c.Request(), err)
		return nil
	}

	if interceptor := frontend.interceptor; interceptor != nil && !interceptor.Intercept(c.Response(), c.Request(), queryValues) {
		return nil
	}

	ServeHTTP(c.Response(), c.Request(), queryValues, frontend.server, frontend.errHandler)
	return nil
}
