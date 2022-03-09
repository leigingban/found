package manger

import (
	"log"

	"github.com/leigingban/found/fundspider"
)

// FetchFundsLatestInfoFrom1234567 更新基金
func (m *Manger) FetchFundsLatestInfoFrom1234567() error {
	// 返回的一个数据中包含多个基金的资料
	dataList, err := fundspider.GetFundInfoByIDsV2(m.FoundCodesListGetter())
	if err != nil {
		log.Println("从网络更新数据时发生错误: ", err)
	}
	for _, data := range dataList {
		found := m.getOrAddFoundByCode(data.FCODE)
		found.UpdateFromData(data)
	}
	return nil
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
