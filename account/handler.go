package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pwlabs/paypal-poc/config"
	"github.com/pwlabs/paypal-poc/oauth"
	"golang.org/x/oauth2"
)

// Handler redirects user to the oauth provider login page
func Handler(env *config.Env) http.Handler {
	oauthConfig := oauth.GetStripeOauthConfig()
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

		w.Write(body)
	})
}
