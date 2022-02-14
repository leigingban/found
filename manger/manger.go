package manger

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/leigingban/found/TTSpider"
	"github.com/leigingban/found/models"
	"github.com/liushuochen/gotable"
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
	table, _ := gotable.Create("代号", "名称", "类型")
	table.Align("名称", gotable.Left)
	table.Align("类型", gotable.Left)
	var fIndex []string
	for key := range m.Founds {
		fIndex = append(fIndex, key)
	}
	sort.Strings(fIndex)
	for _, key := range fIndex {
		fund := m.Founds[key]
		table.AddRow([]string{key, fund.Name, fund.Tags()})
	}
	fmt.Println(table)

	table2, _ := gotable.Create("代号", "名称", "基金数", "类型", "基金")
	table2.Align("名称", gotable.Left)
	table2.Align("类型", gotable.Left)
	table2.Align("基金", gotable.Left)
	var sIndex []string
	for key := range m.Stocks {
		sIndex = append(sIndex, key)
	}
	sort.Strings(sIndex)
	for _, key := range sIndex {
		stock := m.Stocks[key]
		table2.AddRow([]string{key, stock.Name, strconv.Itoa(len(stock.InFunds)), stock.Type, stock.FundNameList()})
	}
	fmt.Println(table2)

}

// ShowInfo 展现数据
func (m *Manger) ShowInfo() {
	// TODO 将代码解耦,排序交给打印的时候操作,已减少后面调整时的代码量
	//table, err := gotable.Create("代号", "名称", "总投入", "最新净值", "最新增量", "预计涨幅", "预计增量")
	table, err := gotable.Create("代号", "名称", "最新净值", "最新增量", "预计涨幅", "预计增量")
	// 0: 居中 1: 左 2:右
	table.Align("名称", 1)
	//table.Align("总投入", 2)
	table.Align("最新净值", 2)
	table.Align("最新增量", 2)
	table.Align("预计涨幅", 2)
	table.Align("预计增量", 2)
	//table.CloseBorder()

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
