package testhelpers

import (
	"net/http"
	"net/http/httptest"
)

//MockEndPoint is used by Stub
type MockEndPoint struct {
	URL     string
	Message []byte
}

//Stub creates a mock http server for testing
func (ep *MockEndPoint) Stub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI != ep.URL {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Write(ep.Message)
		}))
}
