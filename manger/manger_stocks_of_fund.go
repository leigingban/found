package manger

import (
	"github.com/leigingban/found/fundspider"
	"github.com/leigingban/found/models"
)

// NewStockFromRawPtr 创建新Stock
func (m *Manger) NewStockFromRawPtr(rawStock *fundspider.RawStock) *models.Stock {

	/*
		GPDM         string `json:"GPDM"`         // 股票代码
		GPJC         string `json:"GPJC"`         // 股票名称
		JZBL         string `json:"JZBL"`         //
		TEXCH        string `json:"TEXCH"`        //
		ISINVISBL    string `json:"ISINVISBL"`    //
		PCTNVCHGTYPE string `json:"PCTNVCHGTYPE"` //
		PCTNVCHG     string `json:"PCTNVCHG"`     //
		NEWTEXCH     string `json:"NEWTEXCH"`     //
		INDEXCODE    string `json:"INDEXCODE"`    // 股票类别编号
		INDEXNAME    string `json:"INDEXNAME"`    // 股票类别

	*/

	stock := &models.Stock{
		ID:   rawStock.GPDM,
		Name: rawStock.GPJC,
		Type: rawStock.INDEXNAME,
	}

	return stock
}
