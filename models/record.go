package models

import "time"

// Record 记录每一次购买记录，用于统计和计算是否抄底或抄多少
type Record struct {
	TTime time.Time
	Price float64
	Count float64
	value float64
}

// ValueGetter 获取总金额并缓存
func (r *Record) ValueGetter() float64 {
	if r.value == 0 {
		r.value = r.Price * r.Count
	}
	return r.value
}

func (r *Record) LowerThan(other *Record) bool {
	return r.Price < other.Price
}
