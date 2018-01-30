package main

import (
	"net/http"

	"github.com/pwlabs/paypal-poc/handlers"
)

func main() {
	http.Handle("/", handlers.HandleMainPage(`<html><body>
		<a href="/stripeLogin">ConnectWithStripe</a>
		</body></html>`))
	http.HandleFunc("/stripeLogin", handlers.HandleStripeLogin)
	http.HandleFunc("/stripeCallback", handlers.HandleStripeCallback)
	if err := http.ListenAndServe(":65010", nil); err != nil {
		panic(err)
	}
}
