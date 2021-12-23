package models

const BuyPercent uint = 3

type Found struct {
	Id            string
	Name          string
	Amount        float64
	count         float64
	PriceOfMine   float64 //当前基金的平均净值
	PriceEarly    float64 //前一个交易日净值
	PriceRealTime float64 //实时动态净值
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

// AmountGetter 获取总额
func (f *Found) AmountGetter() float64 {
	// 如果有缓存直接返回
	if f.Amount != 0 {
		return f.Amount
	}
	// 累计金额
	var amount float64
	for _, record := range f.Records {
		amount += record.ValueGetter()
	}
	f.Amount = amount
	return f.Amount
}

// CountGetter 获取总得份额
func (f Found) CountGetter() float64 {
	for _, record := range f.Records {
		f.count += record.Count
	}
	return f.count
}

// PriceOfMineGetter 获取等效净值
func (f *Found) PriceOfMineGetter() float64 {
	f.PriceOfMine = f.AmountGetter() / f.CountGetter()
	return f.PriceOfMine
}
