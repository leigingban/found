package httpSpider

import (
	"net/http"
)

var defaultHeader http.Header = map[string][]string{
	"sec-ch-ua":          {"\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\""},
	"sec-ch-ua-mobile":   {"?0"},
	"sec-ch-ua-platform": {"\"Windows\""},
	"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome},96.0.4664.110 Safari/537.36 Edg/96.0.1054.62"},
	"Accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image},webp,image}apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	"Sec-Fetch-Site":     {"none"},
	"Sec-Fetch-Mode":     {"navigate"},
	"Sec-Fetch-User":     {"?1"},
	"Sec-Fetch-Dest":     {"document"},
	"Accept-Language":    {"zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"},
}

// SetHeader 修改自定义的header字段
func SetHeader(key, value string) {
	defaultHeader.Set(key, value)
}
