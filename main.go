package main

import (
	"fmt"
	"github.com/leigingban/found/models"
)

func main() {
	manger := models.CreateManger()
	manger.LoadFromCSV()
	manger.UpToDate()
	fmt.Println(manger)
	fmt.Println(manger)
}
