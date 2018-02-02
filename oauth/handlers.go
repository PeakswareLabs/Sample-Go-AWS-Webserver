package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pwlabs/paypal-poc/config"
	"golang.org/x/oauth2"
)

// MainPageHandler serves the simple html page
func MainPageHandler(html string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, html)
	})
}

// LoginHandler redirects user to the oauth provider login page
func LoginHandler(env *config.Env) http.Handler {
	oauthConfig := env.OauthConfig
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		url := oauthConfig.AuthCodeURL("oauthStateString")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// CallbackHandler gets the oauth token when called by auth provider
func CallbackHandler(env *config.Env) http.Handler {
	oauthConfig := env.OauthConfig
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != "oauthStateString" {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", "oauthStateString", state)
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

		userID := token.Extra("stripe_user_id").(string)
		tokenByte, err := json.Marshal(token)
		err = createOauth(env.DB, &Oauth{AccountID: userID, Token: string(tokenByte)})

		if err != nil {
			fmt.Printf("Database insert failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, `<p>Successfully Authorized Account <code>%s</code>. </p>
		<p>Click <a href="/account?stripe_user_id=%s">here</a> to get account details.</p>
		<p>Click <a href="/oauth/deauthorize?stripe_user_id=%s">here</a> to deauthorize.</p>
		`, userID, userID, userID)
	})
}

// DeauthorizeHandler deauthorizes application with stripe and remvoes oauth token from db
func DeauthorizeHandler(env *config.Env) http.Handler {
	oauthConfig := env.OauthConfig
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stripeUserID := r.FormValue("stripe_user_id")
		if stripeUserID == "" {
			fmt.Printf("invalid stripe user id")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		v := url.Values{
			"stripe_user_id": {stripeUserID},
			"client_id":      {oauthConfig.ClientID},
		}

		_, err := http.NewRequest("POST", "https://connect.stripe.com/oauth/deauthorize", strings.NewReader(v.Encode()))
		if err != nil {
			fmt.Printf("Stripe deauthorize failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		err = deleteOauth(env.DB, stripeUserID)

		if err != nil {
			fmt.Printf("Database delete failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, `<p>Success! Account <code>%s</code> is disconnected.</p>
			<p>Click <a href="{url}">here</a> to restart the OAuth flow.</p>`, stripeUserID)
	})
}
