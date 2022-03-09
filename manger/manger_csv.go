package manger

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/leigingban/found/models"
)

// csvPathGetter 获取文件路径
func (m *Manger) csvPathGetter() string {
	if m.CsvPath != "" {
		return m.CsvPath
	} else {
		m.CsvPath = defaultPath
		return defaultPath
	}
}

// getOrAddFoundByCode 查询或者获取相应的基金
func (m *Manger) getOrAddFoundByCode(foundCode string) *models.Found {
	found, ok := m.Founds[foundCode]
	switch ok {
	case false:
		// 新增基金
		newfound := new(models.Found).New(foundCode)
		m.Founds[foundCode] = newfound
		return newfound
	default:
		return found
	}
}

// NewFundFromLine 通过本地数据创建fund
func (m *Manger) NewFundFromLine(line []string) {

	// 主动检查数据是否满足条件
	if len(line) < 4 {
		log.Println("数据列有缺失，请检查应为: [*,*,*,*] 实际却为: ", line)
		return
	}

	// 每条数据是一个交易记录
	found := m.getOrAddFoundByCode(line[0])

	// 不直接创建是创建fund时需要执行其他步骤
	found.AddRecord(line[1], line[2], line[3])
}

// dataFromCSV 从CSV文件中载入数据
func (m *Manger) dataFromCSV() {

	// 获取数据路径
	dataPath := m.csvPathGetter()

	f, err := os.OpenFile(dataPath, os.O_CREATE, 0777)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	r := csv.NewReader(f)
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		m.NewFundFromLine(line)
	}

}

// SaveToCSV 保存至CSV
func (m *Manger) SaveToCSV() {

	f, err := os.OpenFile(m.csvPathGetter(), os.O_CREATE, 0777)
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
