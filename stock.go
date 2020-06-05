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
