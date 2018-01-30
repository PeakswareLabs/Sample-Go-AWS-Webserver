package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

var oauthConfig *oauth2.Config
var oauthStateString string

func init() {

	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:65010/stripeCallback",
		ClientID:     "ca_CASR50ZBlOjOnrdMaNoav2dcptY7MYx7",
		ClientSecret: "sk_test_Ycr3oc8bxMC4HGYBE5e3ERaY",
		Scopes:       []string{"read_write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://connect.stripe.com/oauth/authorize",
			TokenURL: "https://connect.stripe.com/oauth/token",
		},
	}
	oauthStateString = "random"
}

// HandleMainPage serves the simple html page
func HandleMainPage(html string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html)
	})
}

// HandleStripeLogin redirects user to the stripe oauth login page
func HandleStripeLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleStripeCallback gets the oauth token when called by stripe
func HandleStripeCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, "GOT TOKEN", token.AccessToken)
	// var buffer bytes.Buffer
	// buffer.WriteString("GOT TOKEN ")
	// buffer.WriteString(token.AccessToken)
	// fmt.Printf("GOTT TOKEN == '%s'", token)
	// w.Write(buffer.Bytes())
}
