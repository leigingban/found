package main

import (
	"fmt"
	"github.com/leigingban/found/models"
	"time"
)

func main() {
	r := models.Record{
		TTime: time.Time{},
		Price: 1.8695,
		Count: 2674.51,
	}
	v := r.ValueGetter()
	fmt.Printf("%.2f", v)

	a := models.Found{Id: "008280", Name: "国泰中证煤炭ETF联接C"}

	b := &models.Record{TTime: time.Now(), Price: 1.8695, Count: 2674.51}
	c := &models.Record{TTime: time.Now(), Price: 1.8695, Count: 2674.51}

	a.Records = append(a.Records, b)
	a.Records = append(a.Records, c)

	fmt.Println(a)

	fmt.Printf("%.2f", a.AmountGetter())

}
