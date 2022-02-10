package httpSpider

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

// ProxyOn 打开代理模式
func (hs *HttpSpider) ProxyOn(addrWPort string) {
	tr := hs.client.Transport.(*http.Transport)
	tr.Proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + addrWPort)
	}
}

// DebugOn 打开调试模式，会自动启动 ProxyOn
func (hs *HttpSpider) DebugOn(addrWPort string) {
	hs.ProxyOn(addrWPort)
	tr := hs.client.Transport.(*http.Transport)
	// 取消验证，避免发生错误：x509: certificate signed by unknown authority
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
