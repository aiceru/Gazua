package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tx int

const (
	buy tx = iota
	sell
)

const (
	stockPriceURL = "http://asp1.krx.co.kr/servlet/krx.asp.XMLSiseEng?code="
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
func (s Stock) CalculateIncome(wg *sync.WaitGroup, c chan<- statusReturn, code string) {
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

	c <- statusReturn{code, st}
	wg.Done()
}

func removeCharacter(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}

// getCurrentPrice dummy func
func getCurrentPrice(code string) int {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(stockPriceURL + code)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	var sp StockPrice
	err = xml.Unmarshal(data, &sp)

	curPrice, err := strconv.Atoi(removeCharacter(sp.DailyStocks[0].DayEndPrice, ","))
	if err != nil {
		log.Printf("%v, at getCurrentPrice(code: %v)\n", err, code)
		return 0
	}

	return curPrice
}

// StockPrice represents root element of xml KRX provides
type StockPrice struct {
	XMLName     xml.Name     `xml:"stockprice"`
	DailyStocks []DailyStock `xml:"TBL_DailyStock>DailyStock"`
}

// DailyStock represents table stock info of xml KRX provides
type DailyStock struct {
	DayEndPrice string `xml:"day_EndPrice,attr"`
}
