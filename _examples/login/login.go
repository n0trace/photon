package main

import (
	"bytes"
	"net/url"
	"time"

	"github.com/n0trace/photon/middleware"

	"github.com/n0trace/photon"
)

func main() {
	p := photon.New()
	p.Use(middleware.Limit(time.Second))
	p.Use(middleware.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36"))

	p.Get("https://github.com/login", func(ctx photon.Context) {
		root, err := ctx.Document()
		if err != nil {
			panic(err)
		}
		token, _ := root.Find(`input[name="authenticity_token"]`).Attr("value")

		params := url.Values{}
		params.Add("commit", "Sign in")
		params.Add("authenticity_token", token)
		params.Add("login", "admin@example.com")
		params.Add("password", "password")

		p.Post("https://github.com/session",
			"application/x-www-form-urlencoded",
			bytes.NewBufferString(params.Encode()),
			nil, middleware.FromContext(ctx))
	})
	p.Wait()
}
