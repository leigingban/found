package models

// const BuyPercent 10

type Found struct {
	Id      string
	Name    string
	Amount  float64
	Records []*Record
}

// GetLowestPoint 计算购入的最低点用于后续运算比对
func (f *Found) GetLowestPoint() *Record {
	lower := f.Records[0]
	for i := 1; i < len(f.Records); i++ {
		if f.Records[i].LowerThan(lower) {
			lower = f.Records[i]
		}
	}
	return lower
}

func (f Found) GetValue() {

}

func (f *Found) AmountGetter() float64 {
	for _, record := range f.Records {
		f.Amount += record.ValueGetter()
	}
	return f.Amount
}

func (f Found) GetCount() {

}
