package main

import (
	"net/http"

	"github.com/pwlabs/paypal-poc/oauth"
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
	oauthStateString = "somestring"
}

func main() {

	http.Handle("/", oauth.MainPageHandler(`<html><body>
		<a href="/stripeLogin">ConnectWithStripe</a>
		</body></html>`))
	http.Handle("/stripeLogin", oauth.LoginHandler(oauthStateString, oauthConfig))
	http.Handle("/stripeCallback", oauth.CallbackHandler(oauthStateString, oauthConfig))
	if err := http.ListenAndServe(":65010", nil); err != nil {
		panic(err)
	}
}
