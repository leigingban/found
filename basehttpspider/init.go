package basehttpspider

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// set default timeout
var defaultTimeOut = 3 * time.Second

func newTransport() *http.Transport {
	return &http.Transport{
		Proxy:                  nil, // func(_ *http.Request) (*url.URL, error)
		DialContext:            nil,
		DialTLSContext:         nil,
		TLSClientConfig:        nil, // 可配置是否验证证书 &tls.Config{InsecureSkipVerify: true}
		TLSHandshakeTimeout:    0,
		DisableKeepAlives:      false,
		DisableCompression:     false,
		MaxIdleConns:           10, // 默认是0
		MaxIdleConnsPerHost:    0,
		MaxConnsPerHost:        0,
		IdleConnTimeout:        30 * time.Second, // 默认是0
		ResponseHeaderTimeout:  0,
		ExpectContinueTimeout:  0,
		TLSNextProto:           nil,
		ProxyConnectHeader:     nil,
		GetProxyConnectHeader:  nil,
		MaxResponseHeaderBytes: 0,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}
}

// NewHttpClient 创建一个http客户端，同时初始化内部结构
func NewHttpClient() *http.Client {

	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	// 传输层相关
	client.Transport = newTransport()
	// 设置超时时间
	client.Timeout = defaultTimeOut
	// 设置cookieJar
	client.Jar, _ = cookiejar.New(nil)

	return client
}
