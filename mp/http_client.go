package mp

import (
	"net/http"
	"time"
	"crypto/tls"
	"net"
)

// 一般请求的 http.Client
var TextHttpClient = &http.Client{
	Timeout: 60 * time.Second,
}

// 多媒体上传下载请求的 http.Client
var MediaHttpClient = &http.Client{
	Timeout: 300 * time.Second, // 因为目前微信支持最大的文件是 10MB, 请求超时时间保守设置为 300 秒
}

// NewTLSHttpClient 创建支持双向证书认证的 http.Client
func NewTLSHttpClient(certFile, keyFile string) (httpClient *http.Client, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		},
		Timeout: 60 * time.Second,
	}
	return
}
