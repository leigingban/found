package models

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/leigingban/found/fundspider"
)

const RateToFix float64 = 3

/*
var NowYear int
var NowMonth time.Month
*/

// NowDay 定义一个今日日期
var NowDay int

func init() {
	//NowYear, NowMonth, NowDay = time.Now().Date()
	_, _, NowDay = time.Now().Date()
}

func isToday(dataString string) bool {
	Input, err := time.ParseInLocation("2006-01-02", dataString, Shanghai)
	if err != nil {

		return false
	}
	// 其实比较日就可以了，数据一般在一个月内，不可能出现不同月同日的情况
	//y1, m1, d1 := Input.Date()
	_, _, inputDay := Input.Date()
	//return y1 == NowYear && m1 == NowMonth && d1 == NowDay
	return inputDay == NowDay

}

type Found struct {
	Fundcode      string     // 基金代号
	Name          string     // 基金名称
	DateLatest    *time.Time // 更新(最新)日期
	PriceLatest   float64    // 最新净值 **
	PriceGuess    float64    // 估算净值 **
	RateLatest    float64    // 最新涨幅 **
	RateGuess     float64    // 估算涨幅 **
	Records       []*Record  // 购买记录
	lowestPoint   *Record    // 买入最低点
	Remark        string     // 备注
	notice        string     // 提醒
	LatestIsToday bool       // 最新净值是今天的？
	Stocks        []*Stock
	gc            gcache.Cache
	stockTags     map[string]bool
}

// New 创建一个Found
func (f *Found) New(foundCode string) *Found {
	f.Fundcode = foundCode
	f.Records = []*Record{}
	f.lowestPoint = &Record{}
	f.gc = gcache.New(20).LRU().Build()
	f.stockTags = make(map[string]bool)
	return f
}

// ApplyUpdate 应用获取到的更新数据
func (f *Found) ApplyUpdate(data fundspider.Data) {
	f.Name = data.SHORTNAME
	f.RateGuess = data.GSZZL
	f.PriceGuess = data.GSZ
	f.RateLatest = data.NAVCHGRT
	f.PriceLatest = data.NAV //data.NAV
	// 加入判断最新净值是否今天已更新
	f.LatestIsToday = isToday(data.PDATE)
}

// AddRecord 加入购买记录
func (f *Found) AddRecord(price string, count string, date string) {
	// 添加时将计算好的缓存重设
	record := CreateRecord(price, count, date)

	// 对record进行检查，如果为空，则跳过
	if record == nil {
		return
	}

	f.Records = append(f.Records, record)
}

// iDisEqual 根据id判断Found,用于辅助查询Found
func (f *Found) iDisEqual(fundCode string) bool {
	return f.Fundcode == fundCode
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
	key := "amount"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}

	var amount float64
	for _, record := range f.Records {
		amount += record.LocalBuyAmountGetter()
	}
	f.gc.Set(key, amount)
	return amount
}

// CountGetter 获取总得份额
func (f *Found) CountGetter() float64 {

	key := "countAmount"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}
	// 累计份额
	var count float64
	for _, record := range f.Records {
		count += record.Count
	}
	f.gc.Set(key, count)
	return count
}

// AmountGuessGetter 获取估算总值
func (f *Found) AmountGuessGetter() float64 {
	key := "amountGuest"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}
	amountGuess := f.CountGetter() * f.PriceGuess
	f.gc.Set(key, amountGuess)
	return amountGuess
}

// AmountLatestGetter 获取最新总值
func (f *Found) AmountLatestGetter() float64 {
	key := "amountLatest"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}

	amountLatest := f.CountGetter() * f.PriceLatest
	f.gc.Set(key, amountLatest)
	return amountLatest
}

func (f *Found) PriceBoughtGetter() float64 {

	key := "priceBought"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}

	priceBought := f.AmountBoughtGetter() / f.CountGetter()
	f.gc.Set(key, priceBought)
	return priceBought
}

// AmountRaisedGetter 最新增量
func (f *Found) AmountRaisedGetter() float64 {

	key := "amountRaised"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}

	amountRaised := f.AmountLatestGetter() - f.AmountBoughtGetter()
	f.gc.Set(key, amountRaised)
	return amountRaised
}

// GuestRaisedGetter 预计增量
func (f *Found) GuestRaisedGetter() float64 {

	key := "guestRaised"
	out, err := f.gc.Get(key)
	if err == nil {
		return out.(float64)
	}

	//f.CountGetter()*f.PriceGuess-f.AmountLatestGetter()
	guestRaised := f.CountGetter() * (f.PriceGuess - f.PriceLatest)
	f.gc.Set(key, guestRaised)
	return guestRaised
}

//// Notice 对此基金的提示
//func (f *Found) Notice() string {
//	if f.notice != "" {
//		return f.notice
//	}
//	rateLost := (f.PriceBoughtGetter() - f.PriceGuess) / f.PriceBoughtGetter()
//	if rateLost < RateToFix/100 {
//		return ""
//	}
//	return fmt.Sprintf(" |-*建议: 购入(%.2f)以控制在[%.f%%]\n", f.MoneyToMatchBottom(), RateToFix)
//
//}
//
//// MoneyToMatchBottom 计算保底金额
//func (f *Found) MoneyToMatchBottom() float64 {
//	var money float64
//	moneyLost := (f.PriceBoughtGetter() - f.PriceGuess) * f.CountGetter()
//	totalAmount := 100 * moneyLost / RateToFix
//	money = totalAmount - f.AmountBoughtGetter()
//	return money
//}

func (f *Found) AddStock(stock *Stock) {
	f.Stocks = append(f.Stocks, stock)
	// 将tag作为key用作存储
	f.stockTags[stock.Type] = true
}

func (f *Found) Tags() string {
	var tags []string
	for tag := range f.stockTags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	raw := strings.Builder{}
	l := len(tags)
	for i := 0; i < l; i++ {
		raw.WriteString(tags[i])
		if i < l-1 {
			raw.WriteString(" * ")
		}
	}
	return raw.String()
}

// AmountBoughtStringGetter 总投入
func (f *Found) AmountBoughtStringGetter() string {
	return fmt.Sprintf("%.2f", f.AmountBoughtGetter())
}

// AmountLatestStringGetter 最新净值
func (f *Found) AmountLatestStringGetter() string {
	return fmt.Sprintf("%.2f", f.AmountLatestGetter()/1000)
}

// AmountRaisedStringGetter 最新增量
func (f *Found) AmountRaisedStringGetter() string {
	return fmt.Sprintf("%.2f", f.AmountRaisedGetter()/1000)
}

// CalcTodayRaise 净值当天，反推算出今天的增量 (因为现在只有最新净值，最新涨幅，没其他资料)
func (f *Found) CalcTodayRaise() float64 {
	amount := f.AmountLatestGetter()
	correctRate := f.RateLatest / 100
	return amount * correctRate / (correctRate + 1)
}

// GuestRaisedStringGetter 预计增量
func (f *Found) GuestRaisedStringGetter() string {
	//return fmt.Sprintf("%.2f", f.CountGetter()*f.PriceGuess-f.AmountLatestGetter())

	// 如果最新净值是今天的，预计涨幅已经无意义了
	if f.LatestIsToday {
		// Y = count * PriceLatest /  (RateLatest + 1)
		return fmt.Sprintf("(%.2f)", f.CalcTodayRaise())
	}
	return fmt.Sprintf("%.2f", f.GuestRaisedGetter())
}

// GuestRaisedPercentStringGetter 预计涨幅
func (f *Found) GuestRaisedPercentStringGetter() string {
	// 如果最新净值是今天的，预计涨幅已经无意义了
	if f.LatestIsToday {
		return fmt.Sprintf("(%.2f%%)", f.RateLatest)
	}
	return fmt.Sprintf("%.2f%%", f.RateGuess)
}
