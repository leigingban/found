package main

import (
	"os"

	"github.com/leigingban/found/manger"
)

func main() {
	// 创建一个基金管家
	fundManger := new(manger.Manger).Init()

	// 提取cli后面带的参数
	var args []string
	if len(os.Args) == 1 {
		// 设置一个默认值标识
		args = []string{""}
	} else {
		// 提取参数
		args = os.Args[1:]
	}

	// 根据参数运行相应的代码
	switch args[0] {
	default:
		// 打印即时基金信息
		fundManger.ShowInfo()
	}

	// 获取并打印股票信息
	//fundManger.FetchStocksForFunds()
	//fundManger.PrintStockDetails()

}
