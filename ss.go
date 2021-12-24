package main

import (
	"fmt"
	"github.com/leigingban/found/TTSpider"
)

func main() {
	a, err := TTSpider.GetFundInfoByID("006030")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}
