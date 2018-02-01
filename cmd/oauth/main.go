package main

import (
	"log"
	"net/http"

	"github.com/pwlabs/paypal-poc/account"
	"github.com/pwlabs/paypal-poc/config"
	"github.com/pwlabs/paypal-poc/oauth"
)

func main() {
	db, err := config.NewDB("postgres://labuser:testrdstest@yuva-lambada-test.crw7rg7duol4.us-east-1.rds.amazonaws.com/powerdata")
	if err != nil {
		log.Panic(err)
	}

	env := &config.Env{DB: db}

	http.Handle("/", oauth.MainPageHandler(`<html><body>
		<a href="/oauth/stripeLogin">ConnectWithStripe</a>
		</body></html>`))
	http.Handle("/oauth/stripeLogin", oauth.LoginHandler(env))
	http.Handle("/oauth/stripe_callback", oauth.CallbackHandler(env))
	http.Handle("/oauth/deauthorize", oauth.DeauthorizeHandler(env))
	http.Handle("/account/", account.Handler(env))

	if err := http.ListenAndServe(":65010", nil); err != nil {
		panic(err)
	}
}
