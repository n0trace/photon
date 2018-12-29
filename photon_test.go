package photon_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/n0trace/photon"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(r.Header.Get("User-Agent")))
	})
	return httptest.NewServer(mux)
}

func TestNew(t *testing.T) {
	photon.New()
	photon.New(photon.WithParallel(100))
	ticker := time.NewTicker(time.Second)
	limitFunc := func() <-chan interface{} {
		out := make(chan interface{})
		go func() {
			for t := range ticker.C {
				out <- t
			}
		}()
		return out
	}

	filterFunc := func(r *http.Request) bool {
		return true
	}
	photon.New(photon.WithFilter(filterFunc))
	photon.New(photon.WithLimitFunc(limitFunc))
}
