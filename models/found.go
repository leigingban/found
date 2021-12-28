package models

import (
	"fmt"
	"time"

	"github.com/leigingban/found/TTSpider"
)

const BuyPercent uint = 3

type Found struct {
	Fundcode     string     //基金代号
	Name         string     //基金名称
	DateLarest   *time.Time //更新(最新)日期
	PriceBuy     float64    //买入净值
	PriceLatest  float64    //最新净值
	PriceGuess   float64    //估算净值
	RateLatest   float64    //最新涨幅
	RateGuess    float64    //估算涨幅
	AmountBuy    float64    //买入总值
	AmountLatest float64    //最新总值
	AmountGuess  float64    //估算总值
	Records      []*Record  //购买记录
	lowestPoint  *Record    //买入最低点
	Remark       string     //备注
}

// 创建一个Found
func CreateFound(foundCode string) *Found {
	found := &Found{}
	found.Fundcode = foundCode
	found.Records = []*Record{}
	found.lowestPoint = &Record{}
	return found
}

// GetLowestPoint 计算购入的最低点用于后续运算比对
func (f *Found) GetLowestPoint() *Record {
	// 如果有缓存直接返回
	if f.lowestPoint != nil {
		return f.lowestPoint
	}
	// 计算最低点
	lower := f.Records[0]
	for i := 1; i < len(f.Records); i++ {
		if f.Records[i].LowerThan(lower) {
			lower = f.Records[i]
		}
	}
	// 运算完毕后再赋值
	f.lowestPoint = lower
	return lower
}

// LocalBuyAmountGetter 获取总额
func (f *Found) LocalBuyAmountGetter() float64 {
	// 如果有缓存直接返回
	if f.LocalBuyAmount != 0 {
		return f.LocalBuyAmount
	}
	// 累计金额
	var amount float64
	for _, record := range f.Records {
		amount += record.LocalBuyAmountGetter()
	}
	f.LocalBuyAmount = amount
	return amount
}

// LocalBuyCountGetter 获取总得份额
func (f *Found) LocalBuyCountGetter() float64 {
	// 如果有缓存直接返回
	if f.LocalBuyCount != 0 {
		return f.LocalBuyCount
	}
	// 累计份额
	var count float64
	for _, record := range f.Records {
		count += record.Count
	}
	f.LocalBuyCount = count
	return count
}

// UpdateFromData 从网上更新自身信息
func (f *Found) UpdateFromData(data TTSpider.Data) {
	f.Name = data.SHORTNAME
	f.WebGuessRate = data.GSZZL
	f.WebGuessPrice = data.GSZ
	f.WebFinalRate = data.NAVCHGRT
	f.WebFinalPrice = data.NAV //data.NAV
}

// 获取估算总值
func (f *Found) GuestAmountGetter() float64 {
	if f.GuestAmount != 0 {
		return f.GuestAmount
	}
	f.GuestAmount = f.LocalBuyCountGetter() * f.WebGuessPrice
	return f.GuestAmount
}

// 获取最新总值
func (f *Found) FinalAmountGetter() float64 {
	if f.FinalAmount != 0 {
		return f.FinalAmount
	}
	f.FinalAmount = f.LocalBuyCountGetter() * f.WebFinalPrice
	return f.FinalAmount
}

// 加入购买记录
func (f *Found) AddRecord(price string, count string, date string) {
	// 添加时将计算好的缓存重设
	f.LocalBuyAmount = 0
	f.LocalNowAmount = 0
	f.LocalBuyCount = 0

	record := CreateRecord(price, count, date)

	// 对record进行检查，如果为空，则跳过
	if record == nil {
		return
	}

	f.Records = append(f.Records, record)
}

//根据id判断Found,用于辅助查询Found
func (f *Found) iDisEqual(fundCode string) bool {
	return f.Fundcode == fundCode
}

//计算总投入
func (f Found) LocalBuyAmountToString() string {
	return fmt.Sprintf("%.2f", f.LocalBuyAmount)
}

//展示文本
func (f Found) String() string {
	var raw string
	raw += fmt.Sprintf("--- %s ---\n", f.Name)
	raw += fmt.Sprintf(" |- 代号: %s\n", f.Fundcode)
	raw += fmt.Sprintf(" |- 份额: %.2f\n", f.LocalBuyCountGetter())
	raw += fmt.Sprintf(" |- 预计: %.2f (%.2f)\n", f.GuestAmountGetter(), f.GuestAmountGetter()-f.LocalBuyAmountGetter())
	raw += fmt.Sprintf(" |- 净值: %.2f (%.2f)\n", f.FinalAmountGetter(), f.FinalAmountGetter()-f.LocalBuyAmountGetter())
	raw += fmt.Sprintf(" |- 预涨: %.2f%%\n", f.WebGuessRate)
	raw += fmt.Sprintf(" |- 总涨: %.2f%%\n", (f.FinalAmountGetter()/f.LocalBuyAmountGetter()-1)*100)
	raw += fmt.Sprintf(" |- 备注: %s\n", f.Remark)
	return raw
}
