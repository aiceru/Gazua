package main

import "time"

// User structure for web-user
type User struct {
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_rurl"`
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
