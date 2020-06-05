package main

import (
	"log"
	"time"
)

// User structure for web-user
type User struct {
	UID       string           `json:"uid" bson:"_id,omitempty"`
	Name      string           `json:"name" bson:"name,omitempty"`
	Email     string           `json:"email" bson:"email,omitempty"`
	AvatarURL string           `json:"avatar_url" bson:"avatar_url,omitempty"`
	Stocks    map[string]Stock `json:"stocks" bson:"stocks,omitempty"`
	Expires   time.Time        `json:"expires" bson:"-"`
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
func (u *User) UpdateWithDB(old *User) {
	keys := make([]string, 0)

	if u.Name != old.Name {
		u.Name = old.Name
		keys = append(keys, "name")
	}
	if u.AvatarURL != old.AvatarURL {
		u.AvatarURL = old.AvatarURL
		keys = append(keys, "avatar_url")
	}

	if len(keys) > 0 {
		if err := userdb.UpdateUser(u, keys...); err != nil {
			log.Println("UpdateUser failed: " + err.Error())
		}
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
