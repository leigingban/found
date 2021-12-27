package models

import (
	"fmt"
	"github.com/leigingban/found/TTSpider"
)

const BuyPercent uint = 3

type Found struct {
	Fundcode       string
	Name           string
	WebFinalPrice  float64
	WebGuessPrice  float64
	WebFinalRate   float64
	WebGuessRate   float64
	LocalBuyAmount float64
	LocalNowAmount float64 //当天总额
	LocalBuyCount  float64
	GuestAmount    float64
	FinalAmount    float64
	Records        []*Record
	lowestPoint    *Record
	Remark         string
}

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

func (f *Found) GuestAmountGetter() float64 {
	if f.GuestAmount != 0 {
		return f.GuestAmount
	}
	f.GuestAmount = f.LocalBuyCountGetter() * f.WebGuessPrice
	return f.GuestAmount
}

func (f *Found) FinalAmountGetter() float64 {
	if f.FinalAmount != 0 {
		return f.FinalAmount
	}
	f.FinalAmount = f.LocalBuyCountGetter() * f.WebFinalPrice
	return f.FinalAmount
}

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

func (f *Found) IDisEqual(fundCode string) bool {
	return f.Fundcode == fundCode
}

func (f Found) LocalBuyAmountToString() string {
	return fmt.Sprintf("%.2f", f.LocalBuyAmount)
}

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
