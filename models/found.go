package models

const BuyPercent uint = 3

type Found struct {
	Id            string
	Name          string
	AmountAct     float64 //实际总投入
	AmountNow     float64 //实时净值
	count         float64 //份额
	RateYesterday float64
	RateToday     float64
	RateRealTime  float64
	PriceAct      float64
	PriceToday    float64 //今天的净值
	Records       []*Record
	lowPoint      *Record
}

// GetLowestPoint 计算购入的最低点用于后续运算比对
func (f *Found) GetLowestPoint() *Record {
	// 如果有缓存直接返回
	if f.lowPoint != nil {
		return f.lowPoint
	}
	// 计算最低点
	lower := f.Records[0]
	for i := 1; i < len(f.Records); i++ {
		if f.Records[i].LowerThan(lower) {
			lower = f.Records[i]
		}
	}
	// 运算完毕后再赋值
	f.lowPoint = lower
	return f.lowPoint
}

// AmountActGetter 获取总额
func (f *Found) AmountActGetter() float64 {
	// 如果有缓存直接返回
	if f.AmountAct != 0 {
		return f.AmountAct
	}
	// 累计金额
	var amount float64
	for _, record := range f.Records {
		amount += record.ValueGetter()
	}
	f.AmountAct = amount
	return f.AmountAct
}

// CountGetter 获取总得份额
func (f Found) CountGetter() float64 {
	for _, record := range f.Records {
		f.count += record.Count
	}
	return f.count
}

// PriceActMineGetter 获取等效净值
func (f *Found) PriceActMineGetter() float64 {
	if f.PriceAct != 0 {
		return f.PriceAct
	}
	f.PriceAct = f.AmountActGetter() / f.CountGetter()
	return f.PriceAct
}

// 待完善
func (f *Found) Balance() float64 {
	return 0
}
