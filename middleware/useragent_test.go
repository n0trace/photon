package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/n0trace/photon/middleware"

	"github.com/n0trace/photon"
)

var (
	uajson = `
	{"!Susie":[
        "!Susie (http://www.sync2it.com/susie)"
    ],
    "008":[
        "Mozilla/5.0 (compatible; 008/0.83; http://www.80legs.com/webcrawler.html) Gecko/2008032620"
    ],
    "ABACHOBot":[
        "ABACHOBot"
    ],
    "ABrowse":[
        "Mozilla/5.0 (compatible; U; ABrowse 0.6; Syllable) AppleWebKit/420+ (KHTML, like Gecko)",
        "Mozilla/5.0 (compatible; U; ABrowse 0.6; Syllable) AppleWebKit/420+ (KHTML, like Gecko)",
        "Mozilla/5.0 (compatible; ABrowse 0.4; Syllable)"
		]}
	`
)

func TestRandomUserAgent(t *testing.T) {
	p := photon.New()
	server := newTestServer()
	var useragent string
	var wantUserAgent = "ABACHOBot"
	p.Use(middleware.RandomUserAgent("ABACHOBot"))
	p.OnResponse(func(ctx *photon.Context) {
		useragent, _ = ctx.Response.Text()
	})
	p.Visit(server.URL + "/user-agent")
	p.Wait()
	if useragent != wantUserAgent {
		t.Errorf("RandomUserAgent() useragent = %v, wantUserAgent %v", useragent, wantUserAgent)
		return
	}

	p2 := photon.New()
	p2.Use(middleware.RandomUserAgent("ABACHOBot", "008", "!Susie", "ABrowse"))
	p2.OnResponse(func(ctx *photon.Context) {
		useragent, err := ctx.Response.Text()
		if err != nil {
			t.Errorf("RandomUserAgent() error = %v", err)
		}
		if !strings.Contains(uajson, useragent) {
			t.Errorf("RandomUserAgent() useragent = %v", useragent)
		}
	})

	for i := 0; i < 200; i++ {
		p2.Visit(server.URL+"/user-agent", photon.VisitWithDontFiter())
	}
	p2.Wait()
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(r.Header.Get("User-Agent")))
	})

	mux.HandleFunc("/ua.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(uajson))
	})
	return httptest.NewServer(mux)
}
