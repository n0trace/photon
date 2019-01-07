package middleware_test

import (
	"compress/gzip"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
)

var userMap = map[string]string{
	"hello": "world",
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(r.Header.Get("User-Agent")))
	})

	mux.HandleFunc("/login-cookies", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password, ok := userMap[username]
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if password != r.FormValue("password") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		cookieValue := base64.StdEncoding.EncodeToString([]byte(username))
		cookie := http.Cookie{Name: "userinfo", Value: cookieValue, Path: "/", MaxAge: 86400}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/must-cookies", func(w http.ResponseWriter, r *http.Request) {
		cookies, err := r.Cookie("userinfo")
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		username, err := base64.StdEncoding.DecodeString(cookies.Value)
		_, ok := userMap[string(username)]
		if err != nil || !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})

	mux.HandleFunc("/gzip-hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Flush()
		_, err := gzipWriter.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}
		return
	})

	return httptest.NewServer(mux)
}
