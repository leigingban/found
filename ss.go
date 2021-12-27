package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/leigingban/found/TTSpider"
)

func main() {

	a, _ := TTSpider.GetFundInfoByIDsV2([]string{"005827"})
	spew.Dump(a)
}
