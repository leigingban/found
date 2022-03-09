package manger

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/leigingban/found/models"
	"github.com/liushuochen/gotable"
	"github.com/olekukonko/tablewriter"
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
	// 网上获取数据
	m.update()

	return m
}

// 从网上更新基金信息,作为主要接,从这里调用其他函数
func (m *Manger) update() {
	// 调用本地的更新函数
	err := m.FetchFundsLatestInfoFrom1234567()
	if err != nil {
		log.Fatal("从网上获取数据失败")
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
		{tablewriter.Bold},
		{tablewriter.Bold, tablewriter.FgRedColor},
		{tablewriter.Bold},
		{tablewriter.Bold},
		{tablewriter.Bold, tablewriter.FgRedColor},
		{tablewriter.Bold, tablewriter.FgRedColor},
	}

	if todayRaise <= 0 {
		color[1] = tablewriter.Colors{tablewriter.Bold}
	}
	if GuestRaised <= 0 {
		color[4] = tablewriter.Colors{tablewriter.Bold}
		color[5] = tablewriter.Colors{tablewriter.Bold}
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
