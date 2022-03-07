package manger

// FoundCodesListGetter 获取相应基金列表
func (m Manger) FoundCodesListGetter() []string {
	var raw []string
	for foundCode := range m.Founds {
		raw = append(raw, foundCode)
	}
	return raw
}

// AmountGuessGetter 获取总的估算总值
func (m Manger) AmountGuessGetter() float64 {
	var amount float64
	for _, found := range m.Founds {
		amount += found.AmountGuessGetter()
	}
	return amount
}

// AmountLatestGetter 获取最新的总值
func (m Manger) AmountLatestGetter() float64 {
	var amount float64
	for _, found := range m.Founds {
		amount += found.AmountLatestGetter()
	}
	return amount
}

// AmountBoughtGetter 获取总投入
func (m Manger) AmountBoughtGetter() float64 {
	var amount float64
	for _, found := range m.Founds {
		amount += found.AmountBoughtGetter()
	}
	return amount
}
