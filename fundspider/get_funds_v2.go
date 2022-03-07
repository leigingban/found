package fundspider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func makeFCodes(ids []string) string {
	var raw string
	for _, id := range ids {
		raw += id + "%2C"
	}
	return raw[:len(raw)-3]
}

func GetFundInfoByIDsV2(ids []string) ([]Data, error) {

	url := "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNFInfo"
	method := "POST"

	//ids = "005827%2C003095%2C010806"
	fCodes := makeFCodes(ids)
	payloadRaw := "FCODES=" + fCodes + "&deviceid=D554A7A2-6983-4686-BABD-A5991900AC43&plat=Iphone&product=EFund&version=6.4.9"
	payload := strings.NewReader(payloadRaw)

	client := NewHttpClient()
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return []Data{}, err
	}
	req.Header.Add("User-Agent", "EMProjJijin/6.4.9 (iPad; iOS 15.2; Scale/2.00)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []Data{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []Data{}, err
	}

	raw := new(RawV2)
	err = json.Unmarshal(body, &raw)
	if err != nil {
		log.Println("反序列化时发生错误: ", err)
		//return &Raw{}, err
	}

	return raw.Datas, nil

}
