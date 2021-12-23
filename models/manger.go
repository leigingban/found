package models

type Manger struct {
	Founds      []Found
	totalAmount float64
}

func (m *Manger) TotalAmountGetter() float64 {
	if m.totalAmount != 0 {
		return m.totalAmount
	}
	var total float64
	for _, found := range m.Founds {
		total += found.AmountGetter()
	}
	m.totalAmount = total
	return m.totalAmount
}
