package main

import (
	"bytes"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

const htmlIndex = `<html><body>
<a href="/stripeLogin">ConnectWithStripe</a>
</body></html>`

var (
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
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/stripeLogin", handleStripeLogin)
	http.HandleFunc("/stripeCallback", handleStripeCallback)
	if err := http.ListenAndServe(":65010", nil); err != nil {
		panic(err)
	}
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleStripeLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleStripeCallback(w http.ResponseWriter, r *http.Request) {
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
	var buffer bytes.Buffer
	buffer.WriteString("GOT TOKEN ")
	buffer.WriteString(token.AccessToken)

	fmt.Printf("GOTT TOKEN == '%s'", token)

	w.Write(buffer.Bytes())
}
