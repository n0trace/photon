package middleware_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

func TestFromContext(t *testing.T) {
	p := photon.New()
	p.Get(newTestServer().URL+"/login-cookies", func(ctx photon.Context) {
		switch ctx.StdResponse().StatusCode {
		case http.StatusForbidden:
		default:
			t.Errorf("response status want %v ,got %v", http.StatusForbidden, ctx.StdResponse().StatusCode)
		}
	})

	p.Get(newTestServer().URL+"/must-cookies", func(ctx photon.Context) {
		switch ctx.StdResponse().StatusCode {
		case http.StatusForbidden:
		default:
			t.Errorf("response status want %v ,got %v", http.StatusForbidden, ctx.StdResponse().StatusCode)
		}
	})

	reader := strings.NewReader("username=hello&password=world")

	lastCtx := p.Post(newTestServer().URL+"/login-cookies", "application/x-www-form-urlencoded", reader, func(ctx photon.Context) {
		switch ctx.StdResponse().StatusCode {
		case http.StatusOK:
		default:
			t.Errorf("response status want %v ,got %v", http.StatusOK, ctx.StdResponse().StatusCode)
		}
	})

	p.Get(newTestServer().URL+"/must-cookies", func(ctx photon.Context) {
		switch ctx.StdResponse().StatusCode {
		case http.StatusOK:
		default:
			t.Errorf("response status want %v ,got %v", http.StatusOK, ctx.StdResponse().StatusCode)
		}
	}, middleware.FromContext(lastCtx))
}
