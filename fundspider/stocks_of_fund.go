package fundspider

// 通过基金id获取其下股票的信息

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func fundStocksURLMaker(fundId string) string {
	raw := strings.Builder{}
	raw.WriteString("https://fundmobapi.eastmoney.com/FundMNewApi/FundMNInverstPosition?DATE=")

	// 不确定日期是否要改
	raw.WriteString("2021-12-31")

	raw.WriteString("&FCODE=")

	// 设定基金代码
	raw.WriteString(fundId)

	raw.WriteString("&MobileKey=D554A7A2-6983-4686-BABD-A5991900AC43&OSVersion=15.2.1&appType=ttjj&appVersion=6.5.0&deviceid=D554A7A2-6983-4686-BABD-A5991900AC43&plat=Iphone&product=EFund&serverVersion=6.5.0&version=6.5.0")
	return raw.String()
}

func GetFundStocksByFundId(fundId string) []RawStock {

	url := fundStocksURLMaker(fundId)
	method := "GET"

	client := NewHttpClient()
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return []RawStock{}
	}
	req.Header.Add("Host", "fundmobapi.eastmoney.com")
	req.Header.Add("clientInfo", "ttjj-iPad6,3-iOS-iPadOS15.2.1")
	req.Header.Add("User-Agent", "EMProjJijin/6.5.0 (iPad; iOS 15.2.1; Scale/2.00)")
	req.Header.Add("Referer", "https://mpservice.com/b34ccfc4ed9a4af4a4880fee485cf417/release/pages/fundHold/index")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []RawStock{}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []RawStock{}
	}

	raw := new(RawFundStocks)
	err = json.Unmarshal(body, &raw)
	if err != nil {
		log.Println("反序列化时发生错误: ", err)
		//return &Raw{}, err
	}

	return raw.Datas.FundStocks

}
