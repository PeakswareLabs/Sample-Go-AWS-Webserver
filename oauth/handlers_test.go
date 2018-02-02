package oauth_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pwlabs/paypal-poc/config"
	. "github.com/pwlabs/paypal-poc/oauth"
	"github.com/pwlabs/paypal-poc/testhelpers"
	"golang.org/x/oauth2"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
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

//TODO: Refactor tests and use interface for DB
// Test main for common setup and teardown
// func TestMain(m *testing.M) {
// 	setup()
// 	retCode := m.Run()
// 	tearDown()
// 	os.Exit(retCode)
// }

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
	req := httptest.NewRequest("GET", "/oauth/stripeLogin", nil)
	var oauthConfig = getMockConfig("")
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	env := &config.Env{DB: db, OauthConfig: oauthConfig}

	h := LoginHandler(env)
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

	//setup db mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("create table IF NOT EXISTS oauth").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into oauth").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))

	//Setup stripe mock
	endpoint := testhelpers.MockEndPoint{URL: "/oauth/token", Message: raw}
	server := endpoint.Stub()
	defer server.Close()

	var oauthConfig = getMockConfig(server.URL)
	env := &config.Env{DB: db, OauthConfig: oauthConfig}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/oauth/stripe_callback", nil)
	req.URL.RawQuery = "state=oauthStateString&code=sampleblahcode"
	h := CallbackHandler(env)
	h.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := oauthConfig.Endpoint.TokenURL
	if strings.Contains(req.URL.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
