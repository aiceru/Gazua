package main

import (
	"log"
	"time"
)

// User structure for web-user
type User struct {
	UID       string            `json:"uid" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name,omitempty"`
	Email     string            `json:"email" bson:"email,omitempty"`
	AvatarURL string            `json:"avatar_url" bson:"avatar_url,omitempty"`
	Stocks    map[string]*Stock `json:"-" bson:"stocks,omitempty"`
	Expires   time.Time         `json:"expires" bson:"-"`
}

// Valid check user session validation
func (u *User) Valid() bool {
	return u.Expires.Sub(time.Now()) > 0
}

// Refresh user session
func (u *User) Refresh() {
	u.Expires = time.Now().Add(sessionDuration)
}

// Import googleuser into User struct
func (u *User) Import(g GoogleUser) {
	u.Name = g.Name
	u.Email = g.Email
	u.AvatarURL = g.Picture
}

// UpdateWithDB updates user info, Email is key identifier so that does not change
func (u *User) UpdateWithDB(newUser *User) {
	keys := make([]string, 0)

	if u.Name != newUser.Name {
		u.Name = newUser.Name
		keys = append(keys, "name")
	}
	if u.AvatarURL != newUser.AvatarURL {
		u.AvatarURL = newUser.AvatarURL
		keys = append(keys, "avatar_url")
	}

	if len(keys) > 0 {
		if err := userdb.UpdateUser(u, keys...); err != nil {
			log.Println("UpdateUser failed: " + err.Error())
		}
	}
}

func (u *User) addTx(txType tx, code, name string, quantity, price int) {
	if u.Stocks == nil {
		u.Stocks = make(map[string]*Stock, 0)
	}
	stock, ok := u.Stocks[code]
	if !ok {
		u.Stocks[code] = &Stock{
			Name:  name,
			Buyed: make([]Transaction, 0),
			Sold:  make([]Transaction, 0),
		}
		stock = u.Stocks[code]
	}
	t := Transaction{
		Price:    price,
		Quantity: quantity,
	}

	switch txType {
	case buy:
		stock.Buyed = append(stock.Buyed, t)
	case sell:
		stock.Sold = append(stock.Sold, t)
	default:
		log.Println("Transaction not supported")
	}
}

// GoogleUser represent user info from google oauth2
type GoogleUser struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
	Locale  string `json:"locale"`
}
