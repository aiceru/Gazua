package main

import (
	"net/http"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var renderer *render.Render
var store sessions.Store

func init() {
	// Create new renderer
	renderer = render.New()
}

const (
	sessionKey    = "MBWSserver-session-key"
	sessionSecret = "MBWSserver-session-secret"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/www/*filepath", http.Dir("www"))
	router.GET("/auth/:action/:provider", loginHandler)

	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))
	n.Use(sessionHandler("/auth"))

	n.UseHandler(router)
	n.Run(":9000")
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")

	switch action {
	case "login":
		switch provider {
		case "google":
			loginByGoogle(w, r)
		default:
			http.Error(w,
				"Auth action '"+action+"' with provider '"+provider+"' is not supported.",
				http.StatusNotFound)
		}
	case "callback":
		switch provider {
		case "google":
			authByGoogle(w, r)
		default:
			http.Error(w,
				"Auth action '"+action+"' with provider '"+provider+"' is not supported.",
				http.StatusNotFound)
		}
	default:
		http.Error(w, "Auth action '"+action+"' is not supported.", http.StatusNotFound)
	}
}
