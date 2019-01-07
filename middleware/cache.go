package middleware

import (
	"log"
	"net/http"

	"github.com/n0trace/photon"
)

type CacheConfig struct {
	Driver      CacheDrive
	AllowStatus []int
}

type CacheDrive interface {
	Put(*http.Request, *http.Response) error
}

func CacheWithConfig(config CacheConfig) photon.MiddlewareFunc {
	allowStatusMap := make(map[int]bool)
	for _, status := range config.AllowStatus {
		allowStatusMap[status] = true
	}
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || !ctx.Downloaded() {
				return
			}
			if !allowStatusMap[ctx.StdResponse().StatusCode] {
				return
			}
			err := config.Driver.Put(ctx.Request(), ctx.StdResponse())
			if err != nil {
				log.Println("cache error", err)
			}
		}
	}
}

func Cache() photon.MiddlewareFunc {
	config := CacheConfig{
		Driver:      NewFileDriver(".cache"),
		AllowStatus: []int{http.StatusOK},
	}
	return CacheWithConfig(config)
}
