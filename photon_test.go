package photon_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/n0trace/photon"
)

func TestPhoton(t *testing.T) {
	p := photon.New()
	var executionCount int64
	for i := 0; i < 10; i++ {
		p.Get(newTestServer().URL+"/users?id="+fmt.Sprint(i), func(ctx photon.Context) {
			atomic.AddInt64(&executionCount, 1)
			_, err := ctx.Text()
			if err != nil {
				t.Errorf("ctx.Text().err = %v, want %v", err, nil)
			}
		})
	}
	p.Wait()
	if executionCount != 10 {
		t.Errorf("execution count = %v, want %v", executionCount, 10)
	}
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.RequestURI)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(u.Query().Get("id")))
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{Name: "username", Value: r.FormValue("username"), Path: "/", MaxAge: 86400}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/must-login", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("username")
		w.Write([]byte("hello " + cookie.Value))
	})

	return httptest.NewServer(mux)
}
