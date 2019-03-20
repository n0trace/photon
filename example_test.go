package photon_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

func Example() {
	p := photon.New()
	p.Get(newTestServer().URL+"/users?id=hello", func(ctx photon.Context) {
		text, _ := ctx.Text()
		fmt.Println(text)
	})
	p.Wait()
	//Output:
	//hello
}

func Example_useMiddleware() {
	rootURL := newTestServer().URL
	p := photon.New()
	p.Use(middleware.Limit(200*time.Millisecond), middleware.UserAgent("diy-agent"))
	for i := 0; i != 3; i++ {
		url := fmt.Sprintf("%s/user-agent", rootURL)
		p.Get(url, func(ctx photon.Context) {
			text, _ := ctx.Text()
			fmt.Println(text)
		})
	}
	//or:
	//p.Get(url,callback,middleware...)
	p.Wait()
	//Output:
	//diy-agent
	//diy-agent
	//diy-agent
}

func Example_keepAuth() {
	p := photon.New()

	reader := strings.NewReader("username=foo&password=bar")

	lastCtx := p.Post(newTestServer().URL+"/login",
		"application/x-www-form-urlencoded", reader,
		func(ctx photon.Context) {
			text, _ := ctx.Text()
			fmt.Println(text)
		})

	p.Get(newTestServer().URL+"/must-login", func(ctx photon.Context) {
		text, _ := ctx.Text()
		fmt.Println(text)
	}, middleware.FromContext(lastCtx))

	p.Wait()

	//Output:
	//ok
	//hello foo
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

	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("User-Agent")))
	})

	return httptest.NewServer(mux)
}
