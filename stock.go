package main

// Stock represents stock information
// Stock code (KRX XXXXXX) would be key of the stocks map
type Stock struct {
	Name     string    `json:"name" bson:"name"`
	Holdings []Holding `json:"holdings" bson:"holdings"`
}

// Holding represents stock's price and quantity pair
type Holding struct {
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
	Remain       int     `json:"remain"`
	Income       int     `json:"income"`
	Yield        float32 `json:"yield"`
}

// CalculateIncome calculates current income of stock
func (s Stock) CalculateIncome(code string) *StockStatus {
	st := new(StockStatus)
	st.Name = s.Name
	st.CurrentPrice = getCurrentPrice(code)
	for _, h := range s.Holdings {
		st.Holdings += h.Quantity
		st.Spent += h.Price * h.Quantity
	}
	st.Remain = st.CurrentPrice * st.Holdings
	st.Income = st.Remain - st.Spent
	st.Yield = float32(st.Income) / float32(st.Spent) * 100
	return st
}

// getCurrentPrice dummy func
func getCurrentPrice(code string) int {
	return 15000
}
