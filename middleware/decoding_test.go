package middleware_test

import (
	"strings"
	"testing"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

func TestDecoding(t *testing.T) {
	p := photon.New()
	p.Use(middleware.Header(map[string]string{"Accept-Encoding": "gzip"}))
	p.Use(middleware.Decoding())
	p.Get(newTestServer().URL+"/gzip-hello", func(ctx photon.Context) {
		err := ctx.Error()
		if err != nil {
			panic(err)
		}
		str, _ := ctx.Text()
		if !strings.EqualFold("hello", str) {
			t.Errorf("text want %v ,got %v", "hello", str)
		}
	})
}
