package models

import "github.com/olekukonko/tablewriter"

/*
配置models打印的table
*/

func (f *Found) AddRichRow(table *tablewriter.Table) {
	row := []string{
		f.Fundcode,
		f.Name,
		f.AmountLatestStringGetter(),       // 最新净值
		f.AmountRaisedStringGetter(),       // 最新增量
		f.GuestRaisedPercentStringGetter(), // 预计涨幅
		f.GuestRaisedStringGetter(),        // 预计增量
	}
	color := []tablewriter.Colors{
		{},
		{},
		{},
		{},
		{tablewriter.Bold, tablewriter.FgRedColor},
		{tablewriter.Bold, tablewriter.FgRedColor},
	}
	if f.GuestRaisedGetter() < 0 {
		color[4] = tablewriter.Colors{tablewriter.Bold}
		color[5] = tablewriter.Colors{tablewriter.Bold}
	}
	table.Rich(row, color)
}
