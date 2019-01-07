package middleware_test

import (
	"testing"

	"github.com/n0trace/photon/middleware"

	"github.com/n0trace/photon"
)

func TestRandomUserAgent(t *testing.T) {
	server := newTestServer()
	p := photon.New()
	var expectUserAgent, gotUserAgent string
	expectUserAgent = "diy user-agent"
	p.Use(middleware.UserAgent(expectUserAgent))
	p.Get(server.URL+"/user-agent", func(ctx photon.Context) {
		gotUserAgent, _ = ctx.Text()
	})
	p.Wait()
	if expectUserAgent != gotUserAgent {
		t.Errorf("UserAgent() = %v, want %v", gotUserAgent, expectUserAgent)
	}
}
