package manger

import (
	"fmt"
	"github.com/leigingban/found/TTSpider"
	"github.com/leigingban/found/models"
	"log"
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
}

// ShowInfo 展现数据
func (m *Manger) ShowInfo() {
	fmt.Println(m)
}

//展示文本
func (m Manger) String() string {
	var raw string

	raw += fmt.Sprintf("明细:\n")
	raw += fmt.Sprintf("*总投: %.2f \n", m.AmountBoughtGetter())
	raw += fmt.Sprintf("*预计: %.2f (%.2f)\n", m.AmountGuessGetter(), m.AmountGuessGetter()-m.AmountBoughtGetter())
	raw += fmt.Sprintf("*净值: %.2f (%.2f)\n", m.AmountLatestGetter(), m.AmountLatestGetter()-m.AmountBoughtGetter())

	for _, found := range m.Founds {
		raw += found.String()
	}
	return raw
}
