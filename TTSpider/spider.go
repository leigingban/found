package TTSpider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// BASEURL http://fundgz.1234567.com.cn/js/002190.js?rt=163059678397
const BASEURL = "http://fundgz.1234567.com.cn/js/"

// NewHttpClient 创建一个http客户端
func NewHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}
	return &http.Client{Transport: tr}
}

//CalcSummary 计算两个float并转换成小数点后两位
func CalcSummary(value float64, count float64) float64 {
	raw := value * count
	// 做特别处理 // TODO 待优化
	summary, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", raw), 64)
	return summary
}

/*
Fix1234567Json 对原始数据进行裁切，使其成为标准的json格式
jsonpgz({"fundcode":"001679","name":"前海开源中国稀缺资产混合A","jzrq":"2021-09-02","dwjz":"2.8620","gsz":"2.8302","gszzl":"-1.11","gztime":"2021-09-03 15:00"});
*/
func Fix1234567Json(body []byte) []byte {
	bLen := len(body)
	return body[8 : bLen-2]
}

type Raw struct {
	Fundcode         string  `json:"fundcode"`     // 基金代号
	Name             string  `json:"name"`         // 基金名称
	WebPreviousDate  string  `json:"jzrq"`         // 上一个净值日期
	WebPreviousPrice float64 `json:"dwjz,string"`  // 上一个净值
	WebNowPrice      float64 `json:"gsz,string"`   // 今天净值
	WebNowRate       float64 `json:"gszzl,string"` // 今天增幅
	WebNowTime       string  `json:"gztime"`       // 最后更新时间
}

// GetFundInfoByID 通过爬虫获取原始数据
func GetFundInfoByID(id string) (*Raw, error) {
	// 取前面12位
	rt := strconv.FormatInt(time.Now().UnixMilli(), 10)[0:12]
	url := BASEURL + id + ".js?rt=" + rt
	client := NewHttpClient()

	// 获取数据，遇到错误则返回
	resp, err := client.Get(url)
	if err != nil {
		log.Println("client get 发生错误: ", err)
		return &Raw{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("服务器返回错误，错误如下")
		log.Println("HTTP 错误代号: ", resp.StatusCode)
		log.Println("网址: ", resp.Request.URL)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取http Body时发生错误: ", err)
		return &Raw{}, err
	}

	bodyFixed := Fix1234567Json(body)

	raw := new(Raw)
	err = json.Unmarshal(bodyFixed, &raw)
	if err != nil {
		log.Println("反序列化时发生错误: ", err)
		return &Raw{}, err
	}
	return raw, nil
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
