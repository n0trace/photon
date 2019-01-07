package middleware

import (
	"net/http"

	"github.com/n0trace/photon"
)

type FromCacheDrive interface {
	Get(*http.Request) (*http.Response, bool)
}
type FromCacheConfig struct {
	Driver FromCacheDrive
}

func FromCacheWithConfig(config FromCacheConfig) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			resp, ok := config.Driver.Get(ctx.Request())
			if !ok {
				return
			}
			ctx.SetStdResponse(resp)
			ctx.SetDownload(true)
		}
	}
}

func FromCache() photon.MiddlewareFunc {
	config := FromCacheConfig{
		Driver: NewFileDriver(".cache"),
	}
	return FromCacheWithConfig(config)
}
