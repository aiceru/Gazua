package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Stock represents stock information
// Stock code (KRX XXXXXX) would be key of the stocks map
type Stock struct {
	Name  string        `json:"name" bson:"name"`
	Buyed []Transaction `json:"buyed" bson:"buyed"`
	Sold  []Transaction `json:"sold" bson:"sold"`
}

// Transaction represents stock's price and quantity pair
type Transaction struct {
	Price    int `json:"price" bson:"price"`
	Quantity int `json:"quantity" bson:"quantity"`
}

// StockStatusMap is a map of string - StockStatus pair
type StockStatusMap map[string]*StockStatus

// StockStatus represents current stock's income status
type StockStatus struct {
	Name         string  `json:"name"`
	CurrentPrice int     `json:"current_price"`
	Holdings     int     `json:"holdings"`
	Spent        int     `json:"spent"`
	Earned       int     `json:"earned"`
	Remain       int     `json:"remain"`
	Income       int     `json:"income"` // means remain + earned - spent
	Yield        float32 `json:"yield"`
}

// CalculateIncome calculates current income of stock
func (s Stock) CalculateIncome(code string) *StockStatus {
	st := new(StockStatus)
	st.Name = s.Name
	st.CurrentPrice = getCurrentPrice(code)
	for _, b := range s.Buyed {
		st.Holdings += b.Quantity
		st.Spent += b.Price * b.Quantity
	}
	for _, sold := range s.Sold {
		st.Holdings -= sold.Quantity
		st.Earned += sold.Price * sold.Quantity
	}
	st.Remain = st.CurrentPrice * st.Holdings
	st.Income = st.Remain + st.Earned - st.Spent
	st.Yield = float32(st.Income) / float32(st.Spent) * 100
	return st
}

// getCurrentPrice dummy func
func getCurrentPrice(code string) int {
	return 15000
}

func updateStockCode(name string) string {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://www.google.com/search?q=" + name + "+KRX")
	if err != nil {
		log.Println(err)
		return ""
	}
	code := ""
	z := html.NewTokenizer(resp.Body)
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		switch tt {
		case html.ErrorToken:
			break
		case html.TextToken:
			textToken := string(z.Text())
			if strings.HasSuffix(textToken, "TradingView") {
				fmt.Println(textToken)
				code = strings.TrimLeft(textToken, name + " KRX ")
				code = strings.TrimRight(code, " - TradingView")
				break
			}
		}
	}
	return code
}
