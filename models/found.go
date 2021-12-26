package models

import (
	"fmt"
	"github.com/leigingban/found/TTSpider"
	"log"
)

const BuyPercent uint = 3

type Found struct {
	Fundcode         string
	Name             string
	WebPreviousDate  string
	WebPreviousPrice float64
	WebNowPrice      float64
	WebNowRate       float64
	WebNowTime       string
	LocalBuyAmount   float64
	LocalNowAmount   float64
	LocalBuyCount    float64
	PreviousAmount   float64 // 截止最新总额
	Records          []*Record
	lowestPoint      *Record
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

// PreviousAmountGetter 获取总额
func (f *Found) PreviousAmountGetter() float64 {
	// 如果有缓存直接返回
	if f.PreviousAmount != 0 {
		return f.PreviousAmount
	}
	// 累计金额
	var amount float64
	amount = f.WebPreviousPrice * f.LocalBuyCountGetter()
	f.PreviousAmount = amount
	return amount
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

// UpdateFromWeb 从网上更新自身信息
func (f *Found) UpdateFromWeb() {

	raw, err := TTSpider.GetFundInfoByID(f.Fundcode)
	if err != nil {
		log.Println(err)
	}

	f.Name = raw.Name
	f.WebPreviousDate = raw.WebPreviousDate
	f.WebPreviousPrice = raw.WebPreviousPrice
	f.WebNowPrice = raw.WebNowPrice
	f.WebNowRate = raw.WebNowRate
	f.WebNowTime = raw.WebNowTime

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
	fmt.Println("份额: ", f.LocalBuyCountGetter())
	fmt.Println("净值: ", f.WebPreviousPrice)
	fmt.Println(f.LocalBuyCountGetter() * f.WebPreviousPrice)
	return raw
}
