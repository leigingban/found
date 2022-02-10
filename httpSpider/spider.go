package httpSpider

import (
	"net/http"
	"net/url"
)

const (
	defaultProxyAddr = "http://192.168.4.111:6152"
)

// HttpSpider
// TODO 搞清楚cookie和header的区别
type HttpSpider struct {
	client *http.Client
	header http.Header
	// 添加一个 header
}

// Get 通过在创建request上进行一系列自定义操作，自动挂载默认的header和自定义存储的cookie，header
func (hs *HttpSpider) Get(u string) (*http.Response, error) {
	r, _ := hs.newRequest("GET", u, nil)
	return hs.client.Do(r)
}

// SetCookie 配置cookie
func (hs HttpSpider) SetCookie(u string, session string) {
	nu, _ := url.Parse(u)
	cookies := []*http.Cookie{{
		Name:  "SESSION",
		Value: session},
	}
	hs.client.Jar.SetCookies(nu, cookies)
}

func (hs *HttpSpider) Init() *HttpSpider {
	hs.client = NewHttpClient()
	return hs
}
