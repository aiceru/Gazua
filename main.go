package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var renderer *render.Render
var store sessions.Store
var userdb UserDao

func init() {
	// Create new renderer
	renderer = render.New(render.Options{
		Directory: "web",
	})
	mdao := &MongoDao{
		URL: "mongodb+srv://" + os.Getenv("ATLAS_USER") + ":" +
			os.Getenv("ATLAS_PASS") + "@" + os.Getenv("ATLAS_URI"),
		DBName: databaseName,
	}
	if err := mdao.Connect(); err != nil {
		log.Fatal(err)
	}
	userdb = mdao

	updateCorpList()

	loadCorpMap()
}

const (
	sessionKey    = "MBWSserver-session-key"
	sessionSecret = "MBWSserver-session-secret"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

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
	var sessionUser User
	if session.Get(currentUserKey) != nil {
		// render User info
		data := session.Get(currentUserKey).([]byte)
		json.Unmarshal(data, &sessionUser)
		user, err := userdb.FindUser(sessionUser.Email)
		if err != nil {
			log.Println(err)
			renderer.HTML(w, http.StatusOK, "index", sessionUser)
			return
		}

		stockMap := make(StockStatusMap, 0)
		sum := StockStatus{}
		for code, stock := range user.Stocks {
			income := stock.CalculateIncome(code)
			sum.Spent += income.Spent
			sum.Earned += income.Earned
			sum.Remain += income.Remain
			sum.Income += income.Income
			stockMap[code] = income
		}
		sum.Yield = float32((sum.Income)) / float32(sum.Spent) * 100

		renderer.HTML(w, http.StatusOK, "index", map[string]interface{}{
			"user":     user,
			"stockMap": stockMap,
			"sum":      sum,
		})
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
