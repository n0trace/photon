package middleware

import (
	"time"

	"github.com/n0trace/photon"
)

type LimitFunc func() <-chan interface{}
type LimitConfig struct {
	LimitFunc LimitFunc
}

func LimitWithConfig(config LimitConfig) photon.MiddlewareFunc {
	limitChan := config.LimitFunc()
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			<-limitChan
		}
	}
}

func LimitWithFunc(limitFunc LimitFunc) photon.MiddlewareFunc {
	config := LimitConfig{
		LimitFunc: limitFunc,
	}
	return LimitWithConfig(config)
}

func Limit(duration time.Duration) photon.MiddlewareFunc {
	config := LimitConfig{}
	defaultLimitFunc := func() <-chan interface{} {
		var ticker = time.NewTicker(duration)
		out := make(chan interface{})
		go func() {
			for t := range ticker.C {
				out <- t
			}
		}()
		return out
	}
	if config.LimitFunc == nil {
		config.LimitFunc = defaultLimitFunc
	}
	return LimitWithConfig(config)
}
