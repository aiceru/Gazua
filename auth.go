package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	sessions "github.com/goincremental/negroni-sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	callbackURL         = "http://localhost:9000/auth/callback/google"
	userInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	scopeEmail          = "https://www.googleapis.com/auth/userinfo.email"
	scopeProfile        = "https://www.googleapis.com/auth/userinfo.profile"
)

var oAuthConf *oauth2.Config

// Create random token for "state"
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func loginByGoogle(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Options(sessions.Options{
		Path:   "/auth",
		MaxAge: 300,
	})

	oAuthConf = &oauth2.Config{
		ClientID:     "470304159105-fdeiv96gdmca1i0j1c96gi5jm280nvla.apps.googleusercontent.com",
		ClientSecret: "qFKwFlGkoNMsM6GfXsLGC_e-",
		RedirectURL:  callbackURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{scopeProfile},
	}

	state := randToken()
	session.Set("state", state)
	url := oAuthConf.AuthCodeURL(state)
	log.Printf("Visit the URL for the auth dialog: %v", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func authByGoogle(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	state := session.Get("state")
	session.Delete("state")

	if state != r.FormValue("state") {
		http.Error(w, "Invalid session state", http.StatusUnauthorized)
		return
	}

	token, err := oAuthConf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := oAuthConf.Client(oauth2.NoContext, token)
	userInfoResp, err := client.Get(userInfoAPIEndpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer userInfoResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(string(userInfo))
}

// tok, err := conf.Exchange(oauth2.NoContext, "authorization-code")
// if err != nil {
// 	log.Fatal(err)
// }

// client := conf.Client(oauth2.NoContext, tok)

// // nouse
// fmt.Println(client)
//}
