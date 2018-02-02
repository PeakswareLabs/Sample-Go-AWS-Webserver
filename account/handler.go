package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PeakswareLabs/Go-Webserver/config"
	"github.com/PeakswareLabs/Go-Webserver/oauth"
	"golang.org/x/oauth2"
)

// Handler redirects user to the oauth provider login page
func Handler(env *config.Env) http.Handler {
	oauthConfig := env.OauthConfig
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stripeUserID := r.FormValue("stripe_user_id")
		if stripeUserID == "" {
			fmt.Printf("invalid stripe user id")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		tokenString, err := oauth.RetrieveOauth(env.DB, stripeUserID)

		if err != nil {
			fmt.Printf("Database retrieve failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		token := oauth2.Token{}
		err = json.Unmarshal(tokenString, &token)

		client := oauthConfig.Client(oauth2.NoContext, &token)

		resp, err := client.Get("https://api.stripe.com/v1/accounts/" + stripeUserID)

		if err != nil {
			fmt.Printf("Stripe retrieve failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(body)
	})
}
