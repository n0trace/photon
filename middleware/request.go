package middleware

import (
	"net/http"
	"net/url"
	"time"

	"github.com/n0trace/photon"
)

func Header(headerMap map[string]string) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			request := ctx.Request()
			for k, v := range headerMap {
				request.Header.Set(k, v)
			}
		}
	}
}

func Client(client *http.Client) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			ctx.SetClient(client)
		}
	}
}

func Proxy(u *url.URL) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			ctx.Client().Transport = &http.Transport{
				Proxy: http.ProxyURL(u),
			}
		}
	}
}

func Timeout(t time.Duration) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			ctx.Client().Timeout = t
		}
	}
}

func Transport(transport http.RoundTripper) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			ctx.Client().Transport = transport
		}
	}
}
