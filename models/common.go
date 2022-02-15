package models

import "time"

//Shanghai 标准时区
var Shanghai *time.Location

func init() {
	// BUG FXI "panic: time: missing Location in call to Date"
	// 修复丢失时区错误，在正常window下会有此问题发生，详情直接搜索上述错误

	//var err error
	//Shanghai, err = time.LoadLocation("Asia/Shanghai")
	//if err != nil {
	//	Shanghai = time.FixedZone("CST", 8*60*60)
	//}
	Shanghai = time.FixedZone("CST", 8*60*60)
}
