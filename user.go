package main

import "time"

// User structure for web-user
type User struct {
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Expires   time.Time `json:"expires"`
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
	u.Name = g.FamilyName + " " + g.GivenName
	u.Email = g.Email
	u.AvatarURL = g.Picture
}

// GoogleUser represent user info from google oauth2
type GoogleUser struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	Locale     string `json:"locale"`
}
