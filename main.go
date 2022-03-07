package main

import (
	"github.com/leigingban/found/manger"
)

func main() {
	// 创建一个基金管家
	fundManger := new(manger.Manger).Init()
	// 基金管家从网上获取数据
	fundManger.FetchFundsLatestInfoFrom1234567()
	// 打印即时基金信息
	fundManger.ShowInfo()

	// 获取并打印股票信息
	fundManger.FetchStocksForFunds()
	fundManger.PrintStockDetails()

}
