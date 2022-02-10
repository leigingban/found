package models

type Stock struct {
	ID      string
	Name    string
	Type    string
	InFunds []*Found
}

func (s *Stock) AddFund(fund *Found) {
	s.InFunds = append(s.InFunds, fund)
}
