package TTSpider

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

// BASEURL http://fundgz.1234567.com.cn/js/002190.js?rt=163059678397
const BASEURL = "http://fundgz.1234567.com.cn/js/"

func NewHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}
	return &http.Client{Transport: tr}
}

func GetFundInfoByID(id string) *http.Response {
	// 取前面12位
	rt := strconv.FormatInt(time.Now().UnixMilli(), 10)[0:12]
	url := BASEURL + id + ".js?rt=" + rt
	client := NewHttpClient()

	// 获取数据
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}
	return resp
}

type FoundData struct {
	Fundcode string `json:"fundcode"` // 基金代号
	Name     string `json:"name"`     // 基金名称
	Jzrq     string `json:"jzrq"`     // 上一个净值日期
	Dwjz     string `json:"dwjz"`     // 上一个净值
	Gsz      string `json:"gsz"`      // 今天净值
	Gszzl    string `json:"gszzl"`    // 今天增幅
	Gztime   string `json:"gztime"`   // 最后更新时间
}

/*
返回的数据将是如下格式

jsonpgz({
  "fundcode": "002190",
  "name": "鍐滈摱鏂拌兘婧愪富棰�",
  "jzrq": "2021-12-22",
  "dwjz": "4.3502",
  "gsz": "4.3422",
  "gszzl": "-0.18",
  "gztime": "2021-12-23 15:00"
});

*/
