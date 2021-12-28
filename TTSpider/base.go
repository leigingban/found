package TTSpider

import (
	"net/http"
	"time"
)

// NewHttpClient 创建一个http客户端
func NewHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}
	return &http.Client{Transport: tr}
}

type Data struct {
	FCODE           string      `json:"FCODE"`
	SHORTNAME       string      `json:"SHORTNAME"`
	PDATE           string      `json:"PDATE"`           // 净值日期，一般为上一个交易日净值，如当天是交易日，晚上会更新为当天
	NAV             float64     `json:"NAV,string"`      // 最新净值
	ACCNAV          float64     `json:"ACCNAV,string"`   // 累计净值
	NAVCHGRT        float64     `json:"NAVCHGRT,string"` // 日涨幅
	GSZ             float64     `json:"GSZ,string"`      // 净值估算
	GSZZL           float64     `json:"GSZZL,string"`    // 估算涨幅
	GZTIME          string      `json:"GZTIME"`          // 估算的时间??
	NEWPRICE        interface{} `json:"NEWPRICE"`
	CHANGERATIO     interface{} `json:"CHANGERATIO"`
	ZJL             interface{} `json:"ZJL"`
	HQDATE          interface{} `json:"HQDATE"`
	ISHAVEREDPACKET bool        `json:"ISHAVEREDPACKET"`
}

type Expansion struct {
	GZTIME string `json:"GZTIME"`
	FSRQ   string `json:"FSRQ"`
}

type RawV2 struct {
	Datas        []Data      `json:"Datas"`
	ErrCode      int         `json:"ErrCode"` // 无错误时此项为0
	Success      bool        `json:"Success"`
	ErrMsg       interface{} `json:"ErrMsg"`
	Message      interface{} `json:"Message"`
	ErrorCode    string      `json:"ErrorCode"`
	ErrorMessage interface{} `json:"ErrorMessage"`
	ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
	TotalCount   int         `json:"TotalCount"`
	Expansion    Expansion   `json:"Expansion"`
}
