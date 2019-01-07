package middleware_test

import (
	"strings"
	"testing"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

func TestDecoding(t *testing.T) {
	p := photon.New()
	p.Use(middleware.Decoding())
	p.Get(newTestServer().URL+"/gzip-hello", func(ctx photon.Context) {
		str, _ := ctx.Text()
		if !strings.EqualFold("hello", str) {
			t.Errorf("text want %v ,got %v", "hello", str)
		}
	})
}
