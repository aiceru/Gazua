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

func getSessionUser(r *http.Request) *User {
	s := sessions.GetSession(r)

	currentUserValue := s.Get(currentUserKey)
	if currentUserValue == nil {
		return nil
	}

	var u User
	json.Unmarshal(currentUserValue.([]byte), &u)
	return &u
}

func setSessionUser(r *http.Request, u *User) {
	if u != nil {
		u.Refresh()
	}

	s := sessions.GetSession(r)
	data, _ := json.Marshal(u)
	s.Set(currentUserKey, data)
}

func deleteCurrentUser(r *http.Request) {
	sessions.GetSession(r).Delete(currentUserKey)
}

func sessionHandler(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}

		u := getSessionUser(r)
		if u != nil && u.Valid() {
			setSessionUser(r, u)
			next(w, r)
			return
		}

		deleteCurrentUser(r)
		next(w, r)
		return
	}
}
