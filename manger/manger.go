package manger

import (
	"fmt"
	"github.com/leigingban/found/TTSpider"
	"github.com/leigingban/found/models"
	"github.com/liushuochen/gotable"
	"log"
	"sort"
)

//默认文件路径
const defaultPath = "found.csv"

type Manger struct {
	Founds  map[string]*models.Found //映射，通过id索引相应的found
	Stocks  map[string]*models.Stock
	CsvPath string
}

// Init 初始化一个manger
func (m *Manger) Init() *Manger {
	// 配置其下属性
	m.Founds = make(map[string]*models.Found)
	m.Stocks = make(map[string]*models.Stock)

	// 载入本地数据
	m.dataFromCSV()
	return m
}

// UpdateFundsFromWeb 更新基金
func (m *Manger) UpdateFundsFromWeb() {
	dataList, err := TTSpider.GetFundInfoByIDsV2(m.FoundCodesListGetter())
	if err != nil {
		log.Println("从网络更新数据时发生错误: ", err)
	}
	for _, data := range dataList {
		found := m.getOrAddFoundByCode(data.FCODE)
		found.UpdateFromData(data)
	}
}

// AnalyseFundStocks 分析基金的股票含量
func (m *Manger) AnalyseFundStocks() {
	// 循环遍历本地基金
	for foundId, found := range m.Founds {
		rawStocks := TTSpider.GetFundStocksByFundId(foundId)
		// 循环遍历获取到的数据
		for _, rawStock := range rawStocks {
			// 尝试在本地获取对应的股票
			stock, ok := m.Stocks[rawStock.GPDM]
			// 如果不存在则进行创建
			if !ok {
				stock = m.NewStockFromRawPtr(&rawStock)
				m.Stocks[rawStock.GPDM] = stock
			}
			stock.AddFund(found)
			found.AddStock(stock)
		}
	}
	for _, found := range m.Founds {
		found.Analyse()
	}
}

// ShowInfo 展现数据
func (m *Manger) ShowInfo() {
	//table, err := gotable.Create("代号", "名称", "总投入", "最新净值", "最新增量", "预计涨幅", "预计增量")
	table, err := gotable.Create("代号", "名称", "最新净值", "最新增量", "预计涨幅", "预计增量")
	// 0: 居中 1: 左 2:右
	table.Align("名称", 1)
	//table.Align("总投入", 2)
	table.Align("最新净值", 2)
	table.Align("最新增量", 2)
	table.Align("预计涨幅", 2)
	table.Align("预计增量", 2)

	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return
	}
	// 只是用于排序
	var fundIds []string
	for s := range m.Founds {
		fundIds = append(fundIds, s)
	}
	sort.Strings(fundIds)

	//var AmountBought float64
	var AmountLatest float64
	var AmountRaised float64
	var GuestRaised float64

	for _, id := range fundIds {
		fund := m.Founds[id]
		table.AddRow([]string{
			fund.Fundcode,
			fund.Name,
			//fund.AmountBoughtStringGetter(),
			fund.AmountLatestStringGetter(),
			fund.AmountRaisedStringGetter(),
			fund.GuestRaisedPercentStringGetter(),
			fund.GuestRaisedStringGetter(),
		})
		//AmountBought += fund.AmountBoughtGetter()
		AmountLatest += fund.AmountLatestGetter()
		AmountRaised += fund.AmountRaisedGetter()
		GuestRaised += fund.GuestRaisedGetter()

	}

	table.AddRow([]string{
		"******",
		"合计",
		//fmt.Sprintf("%.2f", AmountBought),
		fmt.Sprintf("%.2f", AmountLatest),
		fmt.Sprintf("%.2f", AmountRaised),
		fmt.Sprintf("%.2f%%", GuestRaised/AmountLatest*100),
		fmt.Sprintf("%.2f", GuestRaised),
	})
	fmt.Println(table)
}
