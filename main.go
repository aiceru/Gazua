package main

import (
	"encoding/json"
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
	renderer = render.New(render.Options{
		Directory: "web",
	})
}

const (
	sessionKey    = "MBWSserver-session-key"
	sessionSecret = "MBWSserver-session-secret"
)

func main() {
	router := httprouter.New()
	router.GET("/", renderMainView)
	router.GET("/auth/:action/:provider", loginHandler)
	router.GET("/logout", logoutHandler)

	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))
	n.Use(sessionHandler("/auth"))

	n.UseHandler(router)
	n.Run(":9000")
}

func renderMainView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessions.GetSession(r)
	var u User
	if session.Get(currentUserKey) != nil {
		// render User info
		data := session.Get(currentUserKey).([]byte)
		json.Unmarshal(data, &u)
		renderer.HTML(w, http.StatusOK, "index", u)
	} else {
		renderer.HTML(w, http.StatusOK, "index", nil)
	}
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

func logoutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sessions.GetSession(r).Delete(currentUserKey)
	http.Redirect(w, r, "/", http.StatusFound)
}
