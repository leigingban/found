package models

import "strings"

type Stock struct {
	ID      string
	Name    string
	Type    string
	InFunds []*Found
}

func (s *Stock) AddFund(fund *Found) {
	s.InFunds = append(s.InFunds, fund)
}

func (s *Stock) FundNameList() string {
	var fundNameList []string
	for _, fund := range s.InFunds {
		fundNameList = append(fundNameList, fund.Name)
	}
	raw := strings.Builder{}
	l := len(fundNameList)
	for i := 0; i < l; i++ {
		raw.WriteString(fundNameList[i])
		if i < l-1 {
			raw.WriteString(" * ")
		}
	}
	return raw.String()
}
