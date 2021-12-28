package models

import (
	"fmt"
	"time"

	"github.com/leigingban/found/TTSpider"
)

const RateToFix float64 = 3

type Found struct {
	Fundcode     string     // 基金代号
	Name         string     // 基金名称
	DateLatest   *time.Time // 更新(最新)日期
	PriceBought  float64    // 买入净值
	PriceLatest  float64    // 最新净值
	PriceGuess   float64    // 估算净值
	RateLatest   float64    // 最新涨幅
	RateGuess    float64    // 估算涨幅
	AmountBought float64    // 买入总值
	AmountLatest float64    // 最新总值
	AmountGuess  float64    // 估算总值
	Count        float64    // 买入数量,份额
	Records      []*Record  // 购买记录
	lowestPoint  *Record    // 买入最低点
	Remark       string     // 备注
	notice       string     // 提醒
}

// CreateFound 创建一个Found
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

// AmountBoughtGetter 获取总额
func (f *Found) AmountBoughtGetter() float64 {
	// 如果有缓存直接返回
	if f.AmountBought != 0 {
		return f.AmountBought
	}
	// 累计金额
	var amount float64
	for _, record := range f.Records {
		amount += record.LocalBuyAmountGetter()
	}
	f.AmountBought = amount
	return amount
}

// CountGetter 获取总得份额
func (f *Found) CountGetter() float64 {
	// 如果有缓存直接返回
	if f.Count != 0 {
		return f.Count
	}
	// 累计份额
	var count float64
	for _, record := range f.Records {
		count += record.Count
	}
	f.Count = count
	return count
}

// UpdateFromData 从网上更新自身信息
func (f *Found) UpdateFromData(data TTSpider.Data) {
	f.Name = data.SHORTNAME
	f.RateGuess = data.GSZZL
	f.PriceGuess = data.GSZ
	f.RateLatest = data.NAVCHGRT
	f.PriceLatest = data.NAV //data.NAV
}

// AmountGuessGetter 获取估算总值
func (f *Found) AmountGuessGetter() float64 {
	if f.AmountGuess != 0 {
		return f.AmountGuess
	}
	f.AmountGuess = f.CountGetter() * f.PriceGuess
	return f.AmountGuess
}

// AmountLatestGetter 获取最新总值
func (f *Found) AmountLatestGetter() float64 {
	if f.AmountLatest != 0 {
		return f.AmountLatest
	}
	f.AmountLatest = f.CountGetter() * f.PriceLatest
	return f.AmountLatest
}

// AddRecord 加入购买记录
func (f *Found) AddRecord(price string, count string, date string) {
	// 添加时将计算好的缓存重设
	f.AmountBought = 0
	f.AmountGuess = 0
	f.AmountLatest = 0
	f.Count = 0

	record := CreateRecord(price, count, date)

	// 对record进行检查，如果为空，则跳过
	if record == nil {
		return
	}

	f.Records = append(f.Records, record)
}

func (f *Found) PriceBoughtGetter() float64 {
	if f.PriceBought != 0 {
		return f.PriceBought
	}
	f.PriceBought = f.AmountBoughtGetter() / f.Count
	return f.PriceBought
}

// iDisEqual 根据id判断Found,用于辅助查询Found
func (f *Found) iDisEqual(fundCode string) bool {
	return f.Fundcode == fundCode
}

// Notice 对此基金的提示
func (f *Found) Notice() string {
	if f.notice != "" {
		return f.notice
	}
	rateLost := (f.PriceBoughtGetter() - f.PriceGuess) / f.PriceBoughtGetter()
	if rateLost < 3/100 {
		return ""
	}
	return fmt.Sprintf(" |-*建议: 购入(%.2f)以控制在[%.f%%]\n", f.MoneyToMatchBottom(), RateToFix)

}

// MoneyToMatchBottom 计算保底金额
func (f *Found) MoneyToMatchBottom() float64 {
	var money float64
	moneyLost := (f.PriceBoughtGetter() - f.PriceGuess) * f.CountGetter()
	totalAmount := 100 * moneyLost / RateToFix
	money = totalAmount - f.AmountBoughtGetter()
	return money
}

//展示文本
func (f Found) String() string {
	var raw string
	raw += fmt.Sprintf("--- %s ---\n", f.Name)
	raw += fmt.Sprintf(" |- 代号: %s\n", f.Fundcode)
	raw += fmt.Sprintf(" |- 份额: %.2f\n", f.CountGetter())
	raw += fmt.Sprintf(" |- 预计: %.2f (%.2f)\n", f.AmountGuessGetter(), f.AmountGuessGetter()-f.AmountBoughtGetter())
	raw += fmt.Sprintf(" |- 净值: %.2f (%.2f)\n", f.AmountLatestGetter(), f.AmountLatestGetter()-f.AmountBoughtGetter())
	raw += fmt.Sprintf(" |- 预涨: %.2f%%\n", f.RateGuess)
	raw += fmt.Sprintf(" |- 总涨: %.2f%%\n", (f.AmountLatestGetter()/f.AmountBoughtGetter()-1)*100)
	raw += fmt.Sprintf(" |- 备注: %s\n", f.Remark)
	raw += f.Notice()
	return raw
}
