package oauth

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// MainPageHandler serves the simple html page
func MainPageHandler(html string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html)
	})
}

// LoginHandler redirects user to the oauth provider login page
func LoginHandler(state string, oauthConfig *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := oauthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// CallbackHandler gets the oauth token when called by auth provider
func CallbackHandler(oauthStateString string, oauthConfig *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}
