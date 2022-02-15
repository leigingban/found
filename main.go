package main

import (
	"github.com/leigingban/found/manger"
)

func main() {
	fundManger := new(manger.Manger).Init()
	fundManger.UpdateFundsFromWeb()
	fundManger.ShowInfo()
	//fundManger.AnalyseFundStocks()

	// 内存分析
	// f, _ := os.OpenFile("mem.profile", os.O_CREATE|os.O_RDWR, 0644)
	// defer f.Close()
	// pprof.Lookup("heap").WriteTo(f, 0)
}
