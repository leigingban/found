package manger

import (
	"fmt"
	"github.com/leigingban/found/fundspider"
	"github.com/leigingban/found/models"
	"github.com/liushuochen/gotable"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"sort"
	"strconv"
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

// FetchFundsLatestInfoFrom1234567 更新基金
func (m *Manger) FetchFundsLatestInfoFrom1234567() {
	// 返回的一个数据中包含多个基金的资料
	dataList, err := fundspider.GetFundInfoByIDsV2(m.FoundCodesListGetter())
	if err != nil {
		log.Println("从网络更新数据时发生错误: ", err)
	}
	for _, data := range dataList {
		found := m.getOrAddFoundByCode(data.FCODE)
		found.UpdateFromData(data)
	}
}

// FetchStocksForFunds 获取基金的股票含量
func (m *Manger) FetchStocksForFunds() {

	// 遍历本地基金，并从网上获取股票信息
	for foundId, found := range m.Founds {
		// 获取当前基金的股票信息
		rawStocks := fundspider.GetFundStocksByFundId(foundId)
		// 遍历获取到的股票，并以键值对形式保存在m.stocks中
		for _, rawStock := range rawStocks {
			// m.stocks中都是唯一值，有则忽略，无则添加
			stock, ok := m.Stocks[rawStock.GPDM]
			if !ok {
				stock = m.NewStockFromRawPtr(&rawStock)
				m.Stocks[rawStock.GPDM] = stock
			}

			// 双向表，各自将数据存放在自身的列表中
			stock.AddFund(found)
			found.AddStock(stock)
		}
	}

}

// PrintStockDetails 打印股票信息
func (m *Manger) PrintStockDetails() {
	/*  创建打印表格
	1. 建立表头
	*/

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

// 生成排好序的funds
func (m *Manger) getSortedFunds() []*models.Found {
	var out []*models.Found

	var fundIds []string
	for s := range m.Founds {
		fundIds = append(fundIds, s)
	}
	sort.Strings(fundIds)

	for _, id := range fundIds {
		out = append(out, m.Founds[id])
	}

	return out
}

func (m *Manger) SetFooter(table *tablewriter.Table) {
	// 统计量
	var AmountLatest float64
	var AmountRaised float64
	var GuestRaised float64
	var todayRaise float64

	for _, fund := range m.Founds {
		// 统计总数，用于footer
		AmountLatest += fund.AmountLatestGetter()
		AmountRaised += fund.AmountRaisedGetter()
		GuestRaised += fund.GuestRaisedGetter()
		if fund.LatestIsToday {
			todayRaise += fund.CalcTodayRaise()
		}
	}

	row := []string{
		"TOTAL:",
		fmt.Sprintf("(%.2f)", todayRaise),
		fmt.Sprintf("%.2f", AmountLatest/1000),
		fmt.Sprintf("%.2f", AmountRaised/1000),
		fmt.Sprintf("%.2f%%", GuestRaised/AmountLatest*100),
		fmt.Sprintf("%.2f", GuestRaised)}

	table.SetFooter(row)

	color := []tablewriter.Colors{
		{},
		{tablewriter.FgRedColor},
		{},
		{},
		{tablewriter.FgRedColor},
		{tablewriter.FgRedColor},
	}

	if todayRaise <= 0 {
		color[1] = tablewriter.Colors{tablewriter.FgGreenColor}
	}
	if GuestRaised <= 0 {
		color[4] = tablewriter.Colors{tablewriter.FgGreenColor}
		color[5] = tablewriter.Colors{tablewriter.FgGreenColor}
	}

	table.SetFooterColor(color...)

}

// ShowInfo 展现数据
func (m *Manger) ShowInfo() {

	table := tablewriter.NewWriter(os.Stdout)

	sortedFunds := m.getSortedFunds()

	for _, fund := range sortedFunds {
		fund.AddRichRow(table)
	}

	table.SetHeader([]string{"编号", "名称", "总k", "增k", "幅", "量"})
	//table.SetFooter(m.getFooterSummary())
	m.SetFooter(table)

	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})

	//table.AppendBulk(data)
	table.Render()
}
