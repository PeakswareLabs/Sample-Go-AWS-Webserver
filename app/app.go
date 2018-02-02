package app

import (
	"log"
	"net/http"

	"github.com/PeakswareLabs/Go-Webserver/account"
	"github.com/PeakswareLabs/Go-Webserver/config"
	"github.com/PeakswareLabs/Go-Webserver/oauth"
)

//Create and initialize all routes
func Create() http.Handler {

	mux := http.NewServeMux()
	db, err := config.NewDB("postgres://labuser:testrdstest@yuva-lambada-test.crw7rg7duol4.us-east-1.rds.amazonaws.com/powerdata")
	if err != nil {
		log.Panic(err)
	}
	oauthConfig := config.GetStripeOauthConfig()

	env := &config.Env{DB: db, OauthConfig: oauthConfig}

	mux.Handle("/", oauth.MainPageHandler())
	mux.Handle("/oauth/stripeLogin", oauth.LoginHandler(env))
	mux.Handle("/oauth/stripe_callback", oauth.CallbackHandler(env))
	mux.Handle("/oauth/deauthorize", oauth.DeauthorizeHandler(env))
	mux.Handle("/account/", account.Handler(env))

	return mux
}
