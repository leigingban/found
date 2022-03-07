package fundspider

// 对基金数据的处理
/*

{
"GPDM": "00700",
"GPJC": "腾讯控股",
"JZBL": "10.12",
"TEXCH": "5",
"ISINVISBL": "0",
"PCTNVCHGTYPE": "增持",
"PCTNVCHG": "0.75",
"NEWTEXCH": "116",
"INDEXCODE": "029026",
"INDEXNAME": "传媒"
}

*/

type RawStock struct {
	GPDM         string `json:"GPDM"`         // 股票代码
	GPJC         string `json:"GPJC"`         // 股票名称
	JZBL         string `json:"JZBL"`         //
	TEXCH        string `json:"TEXCH"`        //
	ISINVISBL    string `json:"ISINVISBL"`    //
	PCTNVCHGTYPE string `json:"PCTNVCHGTYPE"` //
	PCTNVCHG     string `json:"PCTNVCHG"`     //
	NEWTEXCH     string `json:"NEWTEXCH"`     //
	INDEXCODE    string `json:"INDEXCODE"`    // 股票类别编号
	INDEXNAME    string `json:"INDEXNAME"`    // 股票类别
}

type RawFundStocks struct {
	Datas struct {
		FundStocks   []RawStock    `json:"fundStocks"`
		Fundboods    []interface{} `json:"fundboods"`
		Fundfofs     []interface{} `json:"fundfofs"`
		ETFCODE      interface{}   `json:"ETFCODE"`
		ETFSHORTNAME interface{}   `json:"ETFSHORTNAME"`
	} `json:"Datas"`
	ErrCode      int         `json:"ErrCode"`
	Success      bool        `json:"Success"`
	ErrMsg       interface{} `json:"ErrMsg"`
	Message      interface{} `json:"Message"`
	ErrorCode    string      `json:"ErrorCode"`
	ErrorMessage interface{} `json:"ErrorMessage"`
	ErrorMsgLst  interface{} `json:"ErrorMsgLst"`
	TotalCount   int         `json:"TotalCount"`
	Expansion    string      `json:"Expansion"`
}
