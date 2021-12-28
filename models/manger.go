package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/leigingban/found/TTSpider"
)

//默认文件路径
const defaultPath = "found.csv"

type Manger struct {
	Founds  map[string]*Found //映射，通过id索引相应的found
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

// FoundCodesListGetter 获取相应基金列表
func (m Manger) FoundCodesListGetter() []string {
	var raw []string
	for foundCode := range m.Founds {
		raw = append(raw, foundCode)
	}
	return raw
}

// UpToDate 通过爬虫爬取数据并更新found数据
func (m *Manger) UpToDate() {
	dataList, err := TTSpider.GetFundInfoByIDsV2(m.FoundCodesListGetter())
	if err != nil {
		log.Println("从网络更新数据时发生错误: ", err)
	}
	for _, data := range dataList {
		found := m.getOrAddFoundByCode(data.FCODE)
		found.UpdateFromData(data)
	}
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

//展示文本
func (m Manger) String() string {
	var raw string

	raw += fmt.Sprintf("明细:\n")
	raw += fmt.Sprintf("*总投: %.2f \n", m.AmountBoughtGetter())
	raw += fmt.Sprintf("*预计: %.2f (%.2f)\n", m.AmountGuessGetter(), m.AmountGuessGetter()-m.AmountBoughtGetter())
	raw += fmt.Sprintf("*净值: %.2f (%.2f)\n", m.AmountLatestGetter(), m.AmountLatestGetter()-m.AmountBoughtGetter())

	for _, found := range m.Founds {
		raw += found.String()
	}
	return raw
}
