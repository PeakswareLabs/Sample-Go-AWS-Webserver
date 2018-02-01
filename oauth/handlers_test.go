package oauth_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/pwlabs/paypal-poc/oauth"
	"github.com/pwlabs/paypal-poc/testhelpers"
	"golang.org/x/oauth2"
)

func getMockConfig(host string) *oauth2.Config {

	return &oauth2.Config{
		RedirectURL:  "http://localhost/stripeCallback",
		ClientID:     "blah",
		ClientSecret: "shhh",
		Scopes:       []string{"read_write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  host + "/oauth/authorize",
			TokenURL: host + "/oauth/token",
		},
	}
}

//TestMainPageHandler tests the template being served
func TesMainPageHandler(t *testing.T) {
	tables := []string{"hello", "world"}
	for _, sample := range tables {
		h := MainPageHandler(sample)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, nil)
		if status := rec.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if rec.Body.String() != sample {
			t.Errorf("unexpected response: %s", rec.Body.String())
		}
	}
}

//TestLoginHandler tests the login path
func TestLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/StripeLogin", nil)
	var oauthConfig = getMockConfig("")
	h := LoginHandler()
	h.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	//Check the response body is what we expect.
	expected := oauthConfig.Endpoint.AuthURL
	if strings.Contains(req.URL.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}
}

//TestCallbackHandler tests the token path
func TestCallbackHandler(t *testing.T) {
	raw, err := ioutil.ReadFile("fixtures/access_token.json")
	if err != nil {
		t.Error("Unable to open fixture")
	}
	//Setup mock
	endpoint := testhelpers.MockEndPoint{URL: "/oauth/token", Message: raw}
	server := endpoint.Stub()
	defer server.Close()
	var oauthConfig = getMockConfig(server.URL)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/stripeCallback", nil)
	req.URL.RawQuery = "state=random&code=sampleblahcode"
	h := CallbackHandler()
	h.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	expected := oauthConfig.Endpoint.TokenURL
	if strings.Contains(req.URL.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}

}
