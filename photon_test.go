package photon_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(r.Header.Get("User-Agent")))
	})

	mux.HandleFunc("/cookies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			cookies := &http.Cookie{Name: r.FormValue("name"), Value: r.FormValue("value")}
			http.SetCookie(w, cookies)
			w.WriteHeader(http.StatusOK)
		case "GET":
			bs, err := json.Marshal(r.Cookies)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(bs)
		default:
		}
	})

	return httptest.NewServer(mux)
}
