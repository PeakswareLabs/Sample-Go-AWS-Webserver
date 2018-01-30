package main

import (
	"net/http"

	"github.com/pwlabs/paypal-poc/oauth"
)

func main() {
	oauthConfig := oauth.GetStripeOauthConfig()
	http.Handle("/", oauth.MainPageHandler(`<html><body>
		<a href="/stripeLogin">ConnectWithStripe</a>
		</body></html>`))
	http.Handle("/stripeLogin", oauth.LoginHandler("oauthStateString", oauthConfig))
	http.Handle("/stripeCallback", oauth.CallbackHandler("oauthStateString", oauthConfig))
	if err := http.ListenAndServe(":65010", nil); err != nil {
		panic(err)
	}
}
