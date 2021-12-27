package models

import (
	"encoding/csv"
	"fmt"
	"github.com/leigingban/found/TTSpider"
	"io"
	"log"
	"os"
)

const defaultPath = "found.csv"

type Manger struct {
	Founds  map[string]*Found
	CsvPath string
}

// CreateManger 创建一个manger
func CreateManger() *Manger {
	manger := new(Manger)
	manger.LoadFromCSV()
	return manger
}

// CsvPathGetter 获取文件路径
func (m *Manger) CsvPathGetter() string {
	if m.CsvPath != "" {
		return m.CsvPath
	} else {
		m.CsvPath = defaultPath
		return defaultPath
	}
}

// getOrAddFoundByCode 查询或者获取相应的基金
func (m *Manger) getOrAddFoundByCode(foundCode string) *Found {
	found, ok := m.Founds[foundCode]
	switch ok {
	case false:
		newfound := CreateFound(foundCode)
		m.Founds[foundCode] = newfound
		return newfound
	default:
		return found
	}
}

// AddRecord 增加一个记录
func (m *Manger) AddRecord(foundCode string, price string, count string, date string) {
	found := m.getOrAddFoundByCode(foundCode)
	found.AddRecord(price, count, date)
}

// addFromLine 与getOrAddFoundByCode配合使用
func (m *Manger) addFromLine(line []string) {
	if len(line) < 4 {
		log.Println("数据列有缺失，请检查应为: [*,*,*,*] 实际却为: ", line)
	}
	m.AddRecord(line[0], line[1], line[2], line[3])
}

// LoadFromCSV 从CSV文件中载入数据
func (m *Manger) LoadFromCSV() {
	m.Founds = make(map[string]*Found)
	f, err := os.OpenFile(m.CsvPathGetter(), os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		m.addFromLine(line)
	}

}

// SaveToCSV 保存至CSV
func (m *Manger) SaveToCSV() {

	f, err := os.OpenFile(m.CsvPathGetter(), os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	for foundCode, found := range m.Founds {
		for _, record := range found.Records {
			err := w.Write([]string{foundCode, record.PriceToString(), record.CountToString(), record.DateToString()})
			if err != nil {
				log.Println("error happen when writing csv : ", err)
			}
		}
	}
	w.Flush()
	fmt.Println("done")
}

func (m Manger) FoundCodesStringsGetter() []string {
	var raw []string
	for foundCode := range m.Founds {
		raw = append(raw, foundCode)
	}
	return raw
}

func (m *Manger) UpToDate() {
	datas, err := TTSpider.GetFundInfoByIDsV2(m.FoundCodesStringsGetter())
	if err != nil {
		log.Println("从网络更新数据时发生错误: ", err)
	}
	for _, data := range datas {
		found := m.getOrAddFoundByCode(data.FCODE)
		found.UpdateFromData(data)
	}
}

func (m Manger) GuestAmountGetter() float64 {
	var amount float64
	for _, found := range m.Founds {
		amount += found.GuestAmountGetter()
	}
	return amount
}

func (m Manger) FinalAmountGetter() float64 {
	var amount float64
	for _, found := range m.Founds {
		amount += found.FinalAmountGetter()
	}
	return amount
}

func (m Manger) LocalBuyAmountGetter() float64 {
	var amount float64
	for _, found := range m.Founds {
		amount += found.LocalBuyAmountGetter()
	}
	return amount
}

func (m Manger) String() string {
	var raw string

	raw += fmt.Sprintf("明细:\n")
	raw += fmt.Sprintf("*总投: %.2f \n", m.LocalBuyAmountGetter())
	raw += fmt.Sprintf("*预计: %.2f (%.2f)\n", m.GuestAmountGetter(), m.GuestAmountGetter()-m.LocalBuyAmountGetter())
	raw += fmt.Sprintf("*净值: %.2f (%.2f)\n", m.FinalAmountGetter(), m.FinalAmountGetter()-m.LocalBuyAmountGetter())

	for _, found := range m.Founds {
		raw += found.String()
	}
	return raw
}
