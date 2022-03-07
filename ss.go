package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/leigingban/found/fundspider"
)

func main() {

	a, _ := fundspider.GetFundInfoByIDsV2([]string{"005827"})
	spew.Dump(a)
}
