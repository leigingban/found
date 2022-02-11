package main

import (
	"github.com/leigingban/found/manger"
)

func main() {
	fundManger := new(manger.Manger).Init()
	fundManger.UpdateFundsFromWeb()
	fundManger.ShowInfo()
	fundManger.AnalyseFundStocks()
}
