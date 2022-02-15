package models

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

//DATEFORMAT 日期参考格式
//const DATEFORMAT = "2006-01-02 15:04:05"
const DATEFORMAT = "20060102"

// Record 记录每一次购买记录，用于统计和计算是否抄底或抄多少
type Record struct {
	Date       time.Time
	Price      float64
	Count      float64
	value      float64
	nextRecord *Record
}

func CreateRecord(price string, count string, date string) *Record {
	var err error
	record := new(Record)

	record.Price, err = strconv.ParseFloat(price, 64)
	record.Count, err = strconv.ParseFloat(count, 64)
	record.Date, err = time.ParseInLocation(DATEFORMAT, date, Shanghai)

	if err != nil {
		log.Printf("创建Record错误，错误数据: price->%s, count->%s, date->%s", price, count, date)
		log.Println("|-- Err: ", err)
		return nil
	}
	return record
}

// LocalBuyAmountGetter 获取总金额并缓存
func (r *Record) LocalBuyAmountGetter() float64 {
	if r.value == 0 {
		r.value = r.Price * r.Count
	}
	return r.value
}

func (r *Record) LowerThan(other *Record) bool {
	return r.Price < other.Price
}

func (r Record) PriceToString() string {
	return fmt.Sprintf("%.4f", r.Price)
}

func (r Record) CountToString() string {
	return fmt.Sprintf("%.2f", r.Count)
}

func (r Record) DateToString() string {
	return r.Date.Format(DATEFORMAT)
}

func (r Record) String() string {
	var raw string
	raw += fmt.Sprintf("┝-\n")
	return raw
}
