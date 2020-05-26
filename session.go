package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/urfave/negroni"
)

const (
	currentUserKey  = "oauth2_current_user"
	sessionDuration = time.Hour
)

func getCurrentUser(r *http.Request) *User {
	s := sessions.GetSession(r)

	currentUserValue := s.Get(currentUserKey)
	if currentUserValue == nil {
		return nil
	}

	var u User
	json.Unmarshal(currentUserValue.([]byte), &u)
	return &u
}

func setCurrentUser(r *http.Request, u *User) {
	if u != nil {
		u.Refresh()
	}

	s := sessions.GetSession(r)
	data, _ := json.Marshal(u)
	s.Set(currentUserKey, data)
}

func sessionHandler(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}
		next(w, r)
		return
	}
}
