package main

import (
	"fmt"
	"time"
)

// const BuyPercent 10

type Found struct {
	Id      string
	Name    string
	Amount  uint
	Records []Record
}

// 计算购入的最低点用于后续运算比对
func (f *Found) GetLowestPoint() {
	var lowest uint
	for i, r := range f.Records {
		if i == 0 {
			lowest = r.GetValue()
		}
		if value := r.GetValue(); value < lowest {
			lowest = value
		}
	}
	fmt.Println(lowest)
}

func (f Found) GetValue() {

}
func (f Found) GetAmount() {

}
func (f Found) GetCount() {

}

// 记录每一次购买记录，用于统计和计算是否抄底或抄多少
type Record struct {
	TTime time.Time
	Price uint
	Count uint
	value uint
}

func (r *Record) ValueGetter() uint {
	if r.value == 0 {
		r.value = r.Price * r.Count
	}
	return r.value
}

func main() {
	fmt.Println("vim-go")
	a := Record{time.Now(), 1, 2}
	b := Record{time.Now(), 2, 2}
	c := Record{time.Now(), 3, 2}
	d := Record{time.Now(), 4, 2}
	e := Found{"1", "abc", 2, []Record{a, b, c, d}}
	fmt.Println(e)
	e.GetLowestPoint()
}
