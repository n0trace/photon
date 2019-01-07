package middleware

import (
	"sync"

	"github.com/n0trace/photon"
)

type FilterFunc func(photon.Context) bool

type FilterConfig struct {
	FilterFunc FilterFunc
}

func FilterWithConfig(config FilterConfig) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			if ctx.Downloaded() {
				next(ctx)
				return
			}
			filter := config.FilterFunc(ctx)
			if !filter {
				next(ctx)
			}
		}
	}
}

func FilterWithFunc(filterFunc FilterFunc) photon.MiddlewareFunc {
	config := FilterConfig{
		FilterFunc: filterFunc,
	}
	return FilterWithConfig(config)
}

func Filter() photon.MiddlewareFunc {
	config := FilterConfig{}
	crawled := new(sync.Map)
	defaultFilterFunc := func(ctx photon.Context) bool {
		_, ok := crawled.Load(ctx.Request().RequestURI)
		if ok {
			return true
		}
		crawled.Store(ctx.Request().RequestURI, true)
		return false
	}
	if config.FilterFunc == nil {
		config.FilterFunc = defaultFilterFunc
	}
	return FilterWithConfig(config)
}
